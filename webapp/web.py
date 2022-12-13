from pywebio import start_server
from pywebio.input import *
from pywebio.output import *
from pywebio.session import info as session_info
from pywebio.output import put_html
import plotly.express as px
import pandas

def main():
    info = input_group('Basic Parameters', [
        input("Simulation time(s)", name="numGens", type=NUMBER),
        input("Traffic Light (only apply to road with 1 lane)", name="light", type=NUMBER),
        input("Lane Number", name="laneNumber", type=NUMBER),
        input("ratio of SDV/NSDV: lower bound", name="ratioL", type=FLOAT),
        input("ratio of SDV/NSDV: upper bound", name="ratioU", type=FLOAT),
    ])

    a=5
    t=1000

    put_markdown('## Simulation `%d` seconds for a `%d`-lane road without traffic light'%(t,a))
    img = open('D:/cmu_S1/02601/go/src/MultilaneNaschModel/Multilane2.out.gif', 'rb').read()
    put_image(img,width='1000px',height='100px')
    df = pandas.read_csv('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/look.csv')
    fig = px.bar(df, x='Vehicle Type', y='Car flow', color='Vehicle Type',animation_frame="ratio", animation_group="lane", range_y=[0,4000])
    html = fig.to_html(include_plotlyjs="require", full_html=False)
    put_html(html)
    



if __name__ == '__main__':
    #start_server(main, debug=True, port=15200)
    main()