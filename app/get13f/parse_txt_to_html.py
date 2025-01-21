from bs4 import BeautifulSoup
import pandas as pd
import numpy as np
import xml.etree.ElementTree as ET
import requests
import os
import argparse

foot_note = ""

def get_submission(cik):
    # example CIK = 0001649339
    url = f'https://data.sec.gov/submissions/CIK{cik}.json'
    headers = {'User-Agent': 'test-request'}
    response = requests.get(url, headers=headers)
    json_content = response.content

    data = response.json()
    df_filing = pd.DataFrame(data['filings']['recent'])
    return df_filing

def download_filing(cik, accessionNumber):
    # accession number is unique per filing, to remove `-`
    accessionNumber_noDash = accessionNumber.replace('-','')
    
    # local file name is a combination of CIK and accession number
    fname = f'{cik}-{accessionNumber_noDash}'
    local_txt = f'C:/Users/huaji/Docs/13F/{fname}.txt'

    if os.path.exists(local_txt):
        print(f'file already existed at:{local_txt}. no need to download.')
        return
    
    # create url and it's xml table
    url = f'https://www.sec.gov/Archives/edgar/data/{cik}/{accessionNumber_noDash}/{accessionNumber}.txt'
    headers = {'User-Agent': 'Download-SEC-quarter-report'}
    
    # get request to download the data
    response = requests.get(url, headers=headers)

    # write to the file
    with open(local_txt, 'w') as f:
        f.write(response.text)
    
    print(f'File saved at: {local_txt}')


def load_xml(file_path):
    tree = ET.parse(file_path)
    root = tree.getroot()
    return root

def load_txt(file_path):
    # open txt file
    with open(file_path, 'r') as f:
        data = f.read()
    
    # load into xml format
    sp = BeautifulSoup(data, "xml")
    
    # get to the info table tag
    info_tables = sp.find('informationTable')
    xml_string = str(info_tables)
    
    # to xml and return
    xml_element = ET.fromstring(xml_string)
    return xml_element

def extract_infotable(root):
    global foot_note
    data = []
    namespace = {'ns': 'http://www.sec.gov/edgar/document/thirteenf/informationtable'}
    
    for info_table in root.findall('ns:infoTable', namespace):
        entry = {}
        entry['nameOfIssuer'] = info_table.find('ns:nameOfIssuer', namespace).text
        entry['titleOfClass'] = info_table.find('ns:titleOfClass', namespace).text
        entry['cusip'] = info_table.find('ns:cusip', namespace).text
        
        # value and share amount is converted to integers
        entry['value'] = int(info_table.find('ns:value', namespace).text)
        entry['sshPrnamt'] = int(info_table.find('ns:shrsOrPrnAmt/ns:sshPrnamt', namespace).text)
        entry['sshPrnamtType'] = info_table.find('ns:shrsOrPrnAmt/ns:sshPrnamtType', namespace).text
        
        try:
            entry['putCall'] = info_table.find('ns:putCall', namespace).text
        except:
            entry['putCall'] = ''
            
        entry['investmentDiscretion'] = info_table.find('ns:investmentDiscretion', namespace).text

        try:
            entry['otherManager'] = info_table.find('ns:otherManager', namespace).text
        except:
            entry['otherManager'] = ''
        entry['Sole'] = info_table.find('ns:votingAuthority/ns:Sole', namespace).text
        entry['Shared'] = info_table.find('ns:votingAuthority/ns:Shared', namespace).text
        entry['None'] = info_table.find('ns:votingAuthority/ns:None', namespace).text
        data.append(entry)
        df = pd.DataFrame(data)

    if (df['value'] / df['sshPrnamt'] < 1).sum() / df.shape[0] > 0.90:
        foot_note = foot_note + "<p> XML file value are likely in $1,000. </p>"
        df['value'] = df['value']*1000

    df['full-name'] = df['nameOfIssuer'].str.cat(df['titleOfClass'], sep='-')
    return df

def path_builder(cik, accessionNumber):
    accessionNumber_noDash = accessionNumber.replace('-', '')
    url = f'https://www.sec.gov/Archives/edgar/data/{cik}/{accessionNumber_noDash}/{accessionNumber}.txt'
    
    fname = f'{cik}-{accessionNumber_noDash}'
    local_txt = f'C:/Users/huaji/Docs/13F/{fname}.txt'
    return {'url':url, 'local':local_txt}


def merge_pre_post(df_pre, df_post):
    # it will aggregate first
    _df_pre = df_pre.groupby(['full-name','titleOfClass', 'cusip', 'putCall']).agg({'value':'sum', 'sshPrnamt':'sum'}).reset_index()
    _df_curr = df_post.groupby(['full-name','titleOfClass','cusip', 'putCall']).agg({'value':'sum', 'sshPrnamt':'sum'}).reset_index()

    cols = ['full-name','titleOfClass','cusip','value','sshPrnamt','putCall']
    df_all = pd.merge(_df_curr, _df_pre[cols], how='outer', on=['cusip', 'full-name', 'putCall'], suffixes=('_post', '_pre'))
    return df_all

def create_ind(row):
    if pd.isna(row['sshPrnamt_pre']): 
        return 'New'
    elif pd.isna(row['sshPrnamt_post']):
        return 'Sold off'
    elif row['sshPrnamt_post'] == row['sshPrnamt_pre']: 
        return 'No change'
    elif row['sshPrnamt_post'] >= row['sshPrnamt_pre']: 
        return 'Add by {:.1%}'.format(row['sshPrnamt_post']/row['sshPrnamt_pre'] -1)
    else:
        return 'Reduce by {:.1%}'.format(abs(row['sshPrnamt_post']/row['sshPrnamt_pre'] -1))
    


def process_merge_df(df):
    _df = df.copy()
    _df['Change Ind'] = _df.apply(create_ind, axis=1)
#     _df['name'] = np.where(pd.isna(_df['nameOfIssuer_pre']), _df['nameOfIssuer_post'], _df['nameOfIssuer_pre'])
    _df['% of Value'] = _df['value_post']/_df['value_post'].sum()
    _df['avg_price'] = _df['value_post'] / _df['sshPrnamt_post']
    _df['pre_avg_price'] = _df['value_pre'] / _df['sshPrnamt_pre']

    cols_keep = ['full-name', 'putCall', 'Change Ind', 'sshPrnamt_post', 'value_post', 'avg_price', '% of Value', 'sshPrnamt_pre','value_pre','pre_avg_price']
    _df1 = _df[cols_keep]
    _df1.columns = ['Stock', 'Put/Call', 'Change amt', 'Curr. Shares', 'Curr. Value','Curr. AvgPr.', '% of Value', 'Pre. Shares', 'Pre. Value', 'Pre. AvgPr.']
    
    # sort by security group
    grp_by_security = _df1.groupby('Stock').agg({'Curr. Value':'sum'}).sort_values(['Curr. Value'], ascending=False).reset_index()
    grp_by_security.columns = ['Stock', 'total_val']
    
    _result = pd.merge(_df1, grp_by_security, how='inner', on=['Stock'])
    _result = _result.sort_values(['total_val', 'Curr. Value'], ascending=False)
    _result = _result.drop(_result.columns[-1], axis=1)
    
    if (_result['Put/Call'] == '').sum() == _result.shape[0]:
        cols_keep2 = ['Stock', '% of Value', 'Change amt', 'Curr. Shares', 'Curr. Value','Curr. AvgPr.', 'Pre. Shares', 'Pre. Value', 'Pre. AvgPr.']
    else:
        cols_keep2 = ['Stock', '% of Value', 'Change amt', 'Curr. Shares', 'Curr. Value','Curr. AvgPr.', 'Put/Call', 'Pre. Shares', 'Pre. Value', 'Pre. AvgPr.']
    _result = _result[cols_keep2]
    _result = _result.fillna(0)
    _result['Stock'] = _result['Stock'].apply(lambda x: ' '.join([i.capitalize() for i in x.split()]))
    # todo: str.capitalize()
    return _result


def apply_format1(df):
    # apply a pre-defined format for a few columns
    format_dict = {'Curr. Shares': '{:,.0f}', 'Curr. Value': '${:,.0f}', '% of Value':'{:.1%}', 
                   'Curr. AvgPr.': '${:,.2f}', 'Pre. AvgPr.': '${:,.2f}',
                   'Pre. Shares': '{:,.0f}', 'Pre. Value': '${:,.0f}'}
    df_format = df.style.format(format_dict)
    return df_format

def style_negative(column):
    increase = 'color:green;'
    decrease = 'color:red;'
    default = ''
    highlight_lst = []
    for v in column:
        if 'Add' in v or 'New' in v:
            highlight_lst.append(increase)
        elif 'Reduce' in v or 'Sold off' in v:
            highlight_lst.append(decrease)
        else:
            highlight_lst.append(default)
    return highlight_lst

def apply_format(df):
    _df = df.copy()
    _df['Stock'] = _df['Stock'].apply(lambda x: '-'.join(x.split('-')[:-1]))
    
    # apply a pre-defined format for a few columns
    format_dict = {'Curr. Shares': '{:,.0f}', 'Curr. Value': '${:,.0f}', '% of Value':'{:.1%}', 
                   'Curr. AvgPr.': '${:,.1f}', 'Pre. AvgPr.': '${:,.1f}',
                   'Pre. Shares': '{:,.0f}', 'Pre. Value': '${:,.0f}',
                  }
    df_format = _df.style.format(format_dict).apply(style_negative, subset=['Change amt'])                     
    return df_format

style_str = '''
<style>
    table.sortable,
    th {
        border: 1px solid black;
        border-collapse: collapse;
    }

    table.sortable tbody tr:nth-child(odd) {
        background-color: #f5f5f5;
    }

    table.sortable td {
        border-left: 1px solid black;
        border-right: 1px solid black;
    }

    td,
    th {
        padding: 3px 6px 3px 6px;
    }

    th {
        font-size: 85%;
        font-weight: bold;
        text-align: center;
        padding: 4px;
        margin: 1px;
    }

    td {
        font-size: 85%
    }

    .num {
        text-align: right;
    }
</style>
'''

def format_header(header_in_list):
    header_str = "<thead><tr>{}</tr></thead>"\
    .format("".join(["<th>{}</th>".format(cell) for cell in header_in_list]))
    return header_str

def format_body(df):
    format_dict = {'Curr. Shares': '{:,.0f}', 'Curr. Value': '${}', '% of Value':'{:.1%}', 'Put/Call':'{}',
               'Curr. AvgPr.': '${:,.1f}', 'Pre. AvgPr.': '${:,.1f}','Curr. putCall':'{}', 'Pre. putCall':'{}',
               'Pre. Shares': '{:,.0f}', 'Pre. Value': '${}', 'Stock':'{}', 'Change amt':'{}'}
    body_str = "<tbody>"
    row_cnt, col_cnt = df.shape
    for r_i in range(row_cnt):
        row_str ="<tr>"
        row = df.iloc[r_i]
        for c_i in range(col_cnt):
            if row.index[c_i] in ['Stock']:
                row_str += '<td>{}</td>'.format(format_dict[row.index[c_i]].format(row[c_i]))
            elif row.index[c_i] in ['Change amt']:
#                 print(row[c_i],"-----")
#                 print('decrease' in row[c_i].lower(), "----")
                if 'reduce' in row[c_i].lower() or 'sold' in row[c_i].lower():
                    row_str += "<td style=\"color: red\">{}</td>".format(format_dict[row.index[c_i]].format(row[c_i]))
                elif 'add' in row[c_i].lower() or 'new' in row[c_i].lower():
                    row_str += "<td style=\"color: green\">{}</td>".format(format_dict[row.index[c_i]].format(row[c_i]))
                else:
                    row_str += '<td>{}</td>'.format(format_dict[row.index[c_i]].format(row[c_i]))
            else:
                row_str += '<td class="num">{}</td>'.format(format_dict[row.index[c_i]].format(row[c_i]))
        row_str += "</tr>"
        body_str += row_str
        
    body_str += "</tbody>"
    return body_str

def format_number(value):
    # Define magnitude suffixes
    suffixes = ['', 'k', 'm', 'b']

    # Determine the magnitude (thousands, millions, billions)
    magnitude = 0
    while abs(value) >= 1000 and magnitude < len(suffixes) - 1:
        magnitude += 1
        value /= 1000.0

    # Format the value with the appropriate suffix
    formatted_value = f'{value:.1f} {suffixes[magnitude]}'

    return formatted_value

def wrapper(prev, curr):
    # prev = r"C:\Users\huaji\Docs\repos\website-go\db\filings\0001067983\0000950123-23-005270.txt"
    # curr = r"C:\Users\huaji\Docs\repos\website-go\db\filings\0001067983\0000950123-23-008074.txt"

    df_curr = extract_infotable(load_txt(curr))
    df_pre = extract_infotable(load_txt(prev))

    print('<p style="margin: 3px">Current Portfolio: ${:,.0f}</p><p style="margin: 3px">Prior Portfolio: ${:,.0f}</p>'.format(df_curr['value'].sum(), df_pre['value'].sum()))
    df_m = merge_pre_post(df_pre, df_curr)

    result = process_merge_df(df_m)

    result['Curr. Value'] = result['Curr. Value'].apply(lambda x: format_number(x))
    result['Pre. Value'] = result['Pre. Value'].apply(lambda x: format_number(x))

    header = format_header(result.columns)
    body = format_body(result)
    # apply_format(result)
    table_str = "{style}<table class=\"sortable\">{header}{body}</table>".format(style=style_str.replace("\n",""), header=header, body=body)
    print(table_str)
    if foot_note != "":
        print(foot_note)


if __name__ == '__main__':

    argParser = argparse.ArgumentParser()
    argParser.add_argument("-p", "--prev", help="previous file accession number")
    argParser.add_argument("-c", "--curr", help="current file accession number")
    argParser.add_argument("-d", "--dir", help="directory where accession file saved")

    args = argParser.parse_args()

    prev = os.path.join(args.dir, "{}.txt".format(args.prev))
    curr = os.path.join(args.dir, "{}.txt".format(args.curr))

    # print(prev, curr)

    wrapper(prev=prev, curr=curr)

    # prev = r"C:\Users\huaji\Docs\repos\website-go\db\filings\0001067983\0000950123-23-005270.txt"
    # curr = r"C:\Users\huaji\Docs\repos\website-go\db\filings\0001067983\0000950123-23-008074.txt"
   
    # df_curr = extract_infotable(load_txt(curr))
    # df_pre = extract_infotable(load_txt(prev))
    # df_m = merge_pre_post(df_pre, df_curr)

    # print('Curr. Value ${:,.0f}\nPre. Value: ${:,.0f}'.format(df_curr['value'].sum(), df_pre['value'].sum()))

    # result = process_merge_df(df_m)

    # header = format_header(result.columns)
    # body = format_body(result)

    # table_str = "{style}<table class=\"sortable\">{header}{body}</table>".format(style=style_str.replace("\n",""), header=header, body=body)

    # with open(r'C:\Users\huaji\Docs\repos\website-go\test\tmp\html_table.html', 'w') as f:
        # f.write(wrapper(prev=prev, curr=curr))

    

    