from bs4 import BeautifulSoup

# html粘贴到这
html = ''' '''

soup = BeautifulSoup(html)


def get_len(s):
    m = 0
    for c in s:
        if c.isascii():
            m += 1
        else:
            m += 2
    return m


f = True
line_max = []
for r in soup.select('tr'):
    for i, td in enumerate(r.select('td')):
        m = get_len(td.text)
        if not m%2==0:
            m+=1

        if f:
            line_max.append(m)
        else:
            line_max[i] = max(m, line_max[i])
    f=False
sss='+'
for i in line_max:
    sss+='-'*(i+2)
    sss+='+'


with open('table.log','w')as f:
    for r in soup.select('tr'):
        f.write(sss+'\n')
        s='|'
        for i, td in enumerate(r.select('td')):
            t=' '*(line_max[i]-get_len(td.text))
            s+=f' {td.text}{t} |'
        f.write(s+'\n')
    f.write(sss+'\n')

print('表格已输出到table.log')
