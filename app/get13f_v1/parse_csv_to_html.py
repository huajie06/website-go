import pandas as pd

def merge_pre_post(df_pre, df_post):
    # it will aggregate first
    _df_pre = df_pre.groupby(['full-name','TitleOfClass', 'Cusip', 'PutCall']).agg({'Value':'sum', 'SshPrnamt':'sum'}).reset_index()
    _df_curr = df_post.groupby(['full-name','TitleOfClass','Cusip', 'PutCall']).agg({'Value':'sum', 'SshPrnamt':'sum'}).reset_index()

    cols = ['full-name','TitleOfClass','Cusip','Value','SshPrnamt','PutCall']
    df_all = pd.merge(_df_curr, _df_pre[cols], how='outer', on=['Cusip', 'full-name', 'PutCall'], suffixes=('_post', '_pre'))
    return df_all

def create_ind(row):
    if pd.isna(row['SshPrnamt_pre']): 
        return 'New Position'
    elif pd.isna(row['SshPrnamt_post']):
        return 'Sold off'
    elif row['SshPrnamt_post'] == row['SshPrnamt_pre']: 
        return 'No change'
    elif row['SshPrnamt_post'] >= row['SshPrnamt_pre']: 
        return 'Increase by {:.1%}'.format(row['SshPrnamt_post']/row['SshPrnamt_pre'] -1)
    else:
        return 'Decrease by {:.1%}'.format(abs(row['SshPrnamt_post']/row['SshPrnamt_pre'] -1))
    


def process_merge_df(df):
    _df = df.copy()
    _df['Change Ind'] = _df.apply(create_ind, axis=1)
#     _df['name'] = np.where(pd.isna(_df['NameOfIssuer_pre']), _df['NameOfIssuer_post'], _df['NameOfIssuer_pre'])
    _df['% of Value'] = _df['Value_post']/_df['Value_post'].sum()
    _df['avg_price'] = _df['Value_post'] / _df['SshPrnamt_post']
    _df['pre_avg_price'] = _df['Value_pre'] / _df['SshPrnamt_pre']

    cols_keep = ['full-name', 'PutCall', 'Change Ind', 'SshPrnamt_post', 'Value_post', 'avg_price', '% of Value', 'SshPrnamt_pre','Value_pre','pre_avg_price']
    _df1 = _df[cols_keep]
    _df1.columns = ['Stock', 'Curr. PutCall', 'Change amt', 'Curr. Shares', 'Curr. Value','Curr. AvgPr.', '% of Value', 'Pre. Shares', 'Pre. Value', 'Pre. AvgPr.']
    
    # sort by security group
    grp_by_security = _df1.groupby('Stock').agg({'Curr. Value':'sum'}).sort_values(['Curr. Value'], ascending=False).reset_index()
    grp_by_security.columns = ['Stock', 'total_val']
    
    _result = pd.merge(_df1, grp_by_security, how='inner', on=['Stock'])
    _result = _result.sort_values(['total_val', 'Curr. Value'], ascending=False)
    _result = _result.drop(_result.columns[-1], axis=1)
    
    if (_result['Curr. PutCall'] == '').sum() == _result.shape[0]:
        cols_keep2 = ['Stock', '% of Value', 'Change amt', 'Curr. Shares', 'Curr. Value','Curr. AvgPr.', 'Pre. Shares', 'Pre. Value', 'Pre. AvgPr.']
        _result = _result[cols_keep2]
    _result = _result.fillna(0)
    _result['Stock'] = _result['Stock'].apply(lambda x: ' '.join([i.capitalize() for i in x.split()]))
    
    _result.reset_index(drop=True, inplace=True)
    return _result


def process_df(df):
    _df = df.copy()
    if (_df['Value'] / _df['SshPrnamt'] < 1).sum() / _df.shape[0] > 0.90:
        print('share value likely in $1000')
        _df['Value'] = _df['Value']*1000

    _df['full-name'] = _df['NameOfIssuer'].str.cat(df['TitleOfClass'], sep='-')
    
    _df.fillna("", inplace=True)
    return _df

def style_negative(column):
    increase = 'color:green;'
    decrease = 'color:red;'
    default = ''
    highlight_lst = []
    for v in column:
        if 'Increase' in v or 'New Position' in v:
            highlight_lst.append(increase)
        elif 'Decrease' in v or 'Sold off' in v:
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



path = 'C:/Users/huaji/Docs/repos/website-go/app/get13f/filings/0001067983/0000950123-23-008074.csv'
path1 = 'C:/Users/huaji/Docs/repos/website-go/app/get13f/filings/0001067983/0000950123-23-005270.csv'
df = pd.read_csv(path)
df1 = pd.read_csv(path1)

curr = process_df(df1)
pre = process_df(df)

df_all = merge_pre_post(curr, pre)

df_all1 = process_merge_df(df_all)


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

format_dict = {'Curr. Shares': '{:,.0f}', 'Curr. Value': '${:,.0f}', '% of Value':'{:.1%}', 
               'Curr. AvgPr.': '${:,.1f}', 'Pre. AvgPr.': '${:,.1f}',
               'Pre. Shares': '{:,.0f}', 'Pre. Value': '${:,.0f}', 'Stock':'{}', 'Change amt':'{}'}

def format_header(header_in_list):
    header_str = "<thead><tr>{}</tr></thead>"\
    .format("".join(["<th>{}</th>".format(cell) for cell in header_in_list]))
    return header_str

def format_body(df):
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
                if 'decrease' in row[c_i].lower() or 'sold' in row[c_i].lower():
                    row_str += "<td style=\"color: red\">{}</td>".format(format_dict[row.index[c_i]].format(row[c_i]))
                elif 'increase' in row[c_i].lower() or 'new' in row[c_i].lower():
                    row_str += "<td style=\"color: green\">{}</td>".format(format_dict[row.index[c_i]].format(row[c_i]))
                else:
                    row_str += '<td>{}</td>'.format(format_dict[row.index[c_i]].format(row[c_i]))
            else:
                row_str += '<td class="num">{}</td>'.format(format_dict[row.index[c_i]].format(row[c_i]))
        row_str += "</tr>"
        body_str += row_str
        
    
    
    body_str += "</tbody>"
    return body_str


header = format_header(df_all1.columns)
body = format_body(df_all1)


table_str = "{style}<table class=\"sortable\">{header}{body}</table>".format(style=style_str.replace("\n",""), header=header, body=body)


with open('C:/Users/huaji/Docs/repos/website-go/test/tmp/html_table_v2.html', 'w') as f:
    f.write(table_str)