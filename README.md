
# Optimization of traffic capacity based on Cellular Automata

## Instructions

To run the code, enter

`` ./MultilaneNaschModel.exe <sdv_percentage> <simulation_generation> <lane_num> <incident_postion>`` 

Input:

-	sdv_percentage: Ratio of the number of SDV/NSDV generated at the beginning of the road.

-	simulation_generation: Number of generations.

-	lane_num: Number of lanes in the road.

-	incident_position: Position in the road where car crash accidents happen (should be less than 1000).

Output:

-	output.out.gif: Simulation GIF of the road.

## Contributors
  
| Contributors | Content |
| :--:|:--:|
| Ju-Chun Huang| NSDV |
| Xueke Jin| SDV |
| Wanzi Xiao|   SDV |
| Yuxuan Shen| high-level function & website |

## Background
Due to the traffic troubles and the boom in population, urbanization, and economy, the need of transportation grows rapidly. In this project, we suggested different ratio of self-driving vehicles (SDVs) and nonself-driving vehicles(NSDVs) in the road to maximize the traffic capacity.

## Methods
We divided the problem into single-lane and multi-lane scenarios separately in different percentages of SDVs and NSDVs.

- Single-lane

We defined single-lane vehicle acceleration and deceleration based on different human driving situations for NSDVs. SDVs shared the network to make the safety distance smaller than NSDVs, which can provide a larger traffic capacity. Besides, "SDV train" can be an efficient pattern for SDVs.

- Multi-lane

Realistic highway vehicle driving simulation by adding "lane change" situation. SDVs and NSDVs completed turns based on defined rules "Lane-Changing Security(LCS)", "Lane-Changing Motivation(LCM)", and "Lane-Changing Possibility(LCP)".

- Innovation

1. Proposal of SDVs.
2. Addition of traffic lights.
3. Lane change situation addition.

## Results

- Single-lane

1. Changing the percentage of SDVs and NSDVs revealed different traffic capacity outcomes. SDVs formed an SDV train with the same speed in the end.
2. Stopped at traffic lights in an orderly manner, and the SDVs' safe distances are smaller.

- Multi-lane

1. Changing the percentage of SDVs and NSDVs. SDVs formed an SDV train with the same speed in the end.
2. NSDVs turning had a chance due to human factors, and SDVs shared network turns in the most appropriate situation.

- Website

Specified vehicle type, proportions, number of roads, and traffic light placement. Simulated road operation traffic flow.

Since we wrote the website in python, we did not find effective package to run python in GO. We prepared some examples to run in the website. The parameters are listed below:

| Basic Parameters  | eg1  | eg2  | eg3  | eg4  | eg5  | eg6  |
| :---------------: | :--: | ---- | ---- | ---- | ---- | ---- |
| Simulation Times  | 1000 | 1000 | 1000 | 1000 | 1000 | 1000 |
| Incident position |  -1  | -1   | -1   | -1   | 200  | 200  |
|    Lane Number    |  1   | 1    | 5    | 5    | 5    | 5    |
| ratio of SDV/NSDV |  0   | 0.5  | 0.2  | 0.8  | 0.3  | 0.8  |

To open the website, open ```webapp/web.py``` and run the script.

## Video for simulation
[Video](https://drive.google.com/file/d/1MFCtU049tmFCPYqOp1yzHh61lf53hfTW/view?usp=share_link)
