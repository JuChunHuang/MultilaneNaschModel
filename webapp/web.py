from pywebio import start_server
from pywebio.input import *
from pywebio.output import *
from pywebio.session import info as session_info
from pywebio.session import *
from pywebio.output import put_html
import plotly.express as px
import pandas




def main():
    info = input_group('Basic Parameters', [
        input("Simulation time(s)", name="numGens", type=NUMBER),
        input("incident position", name="incipos", type=NUMBER),
        input("Lane Number", name="laneNumber", type=NUMBER),
        input("ratio of SDV/NSDV", name="ratio", type=FLOAT)
    ])

    range0 = input_group("SDV/NSDV ratio range",[
        input("ratio of SDV/NSDV: lower bound", name="ratioL", type=FLOAT),
        input("ratio of SDV/NSDV: upper bound", name="ratioU", type=FLOAT)
    ])



    if info["laneNumber"] == 1:
        if info["ratio"] == 0:
            put_markdown('## Simulation `%d` seconds for a `%d`-lane road'%(1000,1))
            img = open('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/gif/1lane0.out.gif', 'rb').read()
            put_image(img,width='1000px',height='30px')
            df = pandas.read_csv('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/data/1lane.csv')
            fig = px.bar(df, x='Vehicle Type', y='Car flow', color='Vehicle Type',animation_frame="ratio", animation_group="lane", range_y=[0,2000])
            html = fig.to_html(include_plotlyjs="require", full_html=False)
            put_buttons(['Back'], [lambda: go_app('main')])
            put_html(html)
        elif info["ratio"] == 0.5:
            put_markdown('## Simulation `%d` seconds for a `%d`-lane road'%(1000,1))
            img = open('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/gif/1lane50.out.gif', 'rb').read()
            put_image(img,width='1000px',height='30px')
            df = pandas.read_csv('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/data/1lane.csv')
            fig = px.bar(df, x='Vehicle Type', y='Car flow', color='Vehicle Type',animation_frame="ratio", animation_group="lane", range_y=[0,2000])
            html = fig.to_html(include_plotlyjs="require", full_html=False)
            put_buttons(['Back'], [lambda: go_app('main')])
            put_html(html)
        elif info["ratio"] == 1:            
            put_markdown('## Simulation `%d` seconds for a `%d`-lane road'%(1000,1))
            img = open('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/gif/1lane100.out.gif', 'rb').read()
            put_image(img,width='1000px',height='30px')
            df = pandas.read_csv('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/data/1lane.csv')
            fig = px.bar(df, x='Vehicle Type', y='Car flow', color='Vehicle Type',animation_frame="ratio", animation_group="lane", range_y=[0,2000])
            html = fig.to_html(include_plotlyjs="require", full_html=False)
            put_buttons(['Back'], [lambda: go_app('main')])
            put_html(html)

    if info["laneNumber"] == 5 and info["incipos"] == -1:
        if info["ratio"] == 0.2:
            put_markdown('## Simulation `%d` seconds for a `%d`-lane road'%(1000,5))
            img = open('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/gif/5lane20sdv.out.gif', 'rb').read()
            put_image(img,width='1000px',height='100px')
            df = pandas.read_csv('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/data/5lane.csv')
            fig = px.bar(df, x='Vehicle Type', y='Car flow', color='Vehicle Type',animation_frame="ratio", animation_group="lane", range_y=[0,5000])
            html = fig.to_html(include_plotlyjs="require", full_html=False)
            put_buttons(['Back'], [lambda: go_app('main')])
            put_html(html)
        elif info["ratio"] == 0.8:
            put_markdown('## Simulation `%d` seconds for a `%d`-lane road'%(1000,5))
            img = open('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/gif/5lane80sdv.out.gif', 'rb').read()
            put_image(img,width='1000px',height='100px')
            df = pandas.read_csv('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/data/5lane.csv')
            fig = px.bar(df, x='Vehicle Type', y='Car flow', color='Vehicle Type',animation_frame="ratio", animation_group="lane", range_y=[0,5000])
            html = fig.to_html(include_plotlyjs="require", full_html=False)
            put_buttons(['Back'], [lambda: go_app('main')])
            put_html(html)
    if info["laneNumber"] == 5 and info["incipos"] == 200:
        if info["ratio"] == 0.3:
            put_markdown('## Simulation `%d` seconds for a `%d`-lane road'%(1000,5))
            img = open('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/gif/mid30sdv.out.gif', 'rb').read()
            put_image(img,width='1000px',height='100px')
            df = pandas.read_csv('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/data/5laneInci.csv')
            fig = px.bar(df, x='Vehicle Type', y='Car flow', color='Vehicle Type',animation_frame="ratio", animation_group="lane", range_y=[0,5000])
            html = fig.to_html(include_plotlyjs="require", full_html=False)
            put_buttons(['Back'], [lambda: go_app('main')])
            put_html(html)
        elif info["ratio"] == 0.8:
            put_markdown('## Simulation `%d` seconds for a `%d`-lane road'%(1000,5))
            img = open('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/gif/mid80sdv.out.gif', 'rb').read()
            put_image(img,width='1000px',height='100px')
            df = pandas.read_csv('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/data/5laneInci.csv')
            fig = px.bar(df, x='Vehicle Type', y='Car flow', color='Vehicle Type',animation_frame="ratio", animation_group="lane", range_y=[0,5000])
            html = fig.to_html(include_plotlyjs="require", full_html=False)
            put_buttons(['Back'], [lambda: go_app('main')])
            put_html(html)
    else:
        put_markdown('## Simulation `%d` seconds for a `%d`-lane road'%(1000,5))
        img = open('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/gif/mid80sdv.out.gif', 'rb').read()
        put_image(img,width='1000px',height='100px')
        df = pandas.read_csv('D:/cmu_S1/02601/go/src/MultilaneNaschModel/webapp/data/5laneInci.csv')
        fig = px.bar(df, x='Vehicle Type', y='Car flow', color='Vehicle Type',animation_frame="ratio", animation_group="lane", range_y=[0,5000])
        html = fig.to_html(include_plotlyjs="require", full_html=False)
        put_buttons(['Back'], [lambda: go_app('main')])
        put_html(html)
        
            

    
    



if __name__ == '__main__':
    start_server(main, debug=True, port=15200)
    #main()