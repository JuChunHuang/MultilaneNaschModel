
<div align='center' ><font size='6'>Optimization of traffic capacity based on Cellular Automata</font></div>

- [Instructions](#Instructions)
- [Contributor](#Contributor)
- [Background](#Background)
- [Method](#Method)
  - [Siglelane](#Siglelane])
  - [Multilane](#Multilane)
  - [Innovation](#Innovation)
- [Results](#Results)
  - [Siglelane](#Siglelane])
  - [Multilane](#Multilane)
  - [Website](#Webset)
- [Expectation](#Expectation)
- [Video](#Video)



* # Instructions

To run the code, enter`` .\MultilaneNaschModel.exe <sdv_percentage> <simulation_generation> <laneNum> <incident_postion>`` 

Input:

​	sdv_percentage: The ratio of SDV/NSDV generated at the beginning of the road

​	simulation_generation: The number of generations

​	laneNum: The number of lanes in the road

​	incident_position: the position in the road where car crash accidents happen，should be less than 1000.

output:

​	output.out.gif: the simulation GIF of the road

* # Contributor
  
| Contributor | Content |
| :--:|:--:|
| Ju-Chun Huang| NSDV |
| Xueke Jin| SDV |
| Wanzi Xiao|   SDV |
| Yuxuan Shen| highlevel function&website |

* # Background
Due to more traffic troubles due to the boom in population, urbanization, and economy, the need for transportation grows rapidly. In this project, we suggest self-driving vehicles(SDVs) replace nonself-driving vehicles(NSDVs) to optimize traffic capacity.

* # Method
We divide the problem into single-lane and multi-lane separately in different percentages of SDVs and NSDVs. Check running situations on SDVs and NSDVs and traffic flow on the lane.
  ## Singlelane
Define for single lane vehicle acceleration and deceleration, NSDVs have human driving situation. SDVs can share the network to make the safety distance smaller to get more traffic flow.Besides SDV train canbe an efficient pattern for SDVs.
  ## Multilane
Realistic highway vehicle driving simulation by adding lane change content. SDVs and NSDVs complete turns with defined rules Lane-Changing Security(LCS), Lane-Changing Motivation(LCM), and Lane-Changing Possibility(LCP).
 ## Innovation
1. Proposal of SDVs
2. Addition of traffic lights
3. Lane change situation addition
* # Results
  ## Singlelane
1. Change the percentage of SDVs and NSDVs. SDVs are forming an SDV train in the end and go at the same speed as all the other vehicles in the train and have an obvious pattern.
2. Stop at traffic lights in an orderly manner, and the SDVs' following spacing is smaller.
  ## Multilane
1. Change the percentage of SDVs and NSDVs. SDVs are forming an SDV train in the end and go at the same speed as all the other vehicles in the train and have an obvious pattern.
2. NSDVs Turning has a chance due to human factors, and SDVs share network turns in the most appropriate situations.
  ## Website
Specify vehicle type, proportion, number of roads, and traffic light placement. Simulate road operation traffic flow.

Since we wrote the website in python, and did not find effective package to run python in go, we prepare some example to run in the website. The parameters are listed below:

| Basic Parameters  | eg1  | eg2  | eg3  | eg4  | eg5  | eg6  |
| :---------------: | :--: | ---- | ---- | ---- | ---- | ---- |
| Simulation Times  | 1000 | 1000 | 1000 | 1000 | 1000 | 1000 |
| Incident position |  -1  | -1   | -1   | -1   | 200  | 200  |
|    Lane Number    |  1   | 1    | 5    | 5    | 5    | 5    |
| ratio of SDV/NSDV |  0   | 0.5  | 0.2  | 0.8  | 0.3  | 0.8  |

To run the website, open ```webapp/web.py``` and run the script.

* # Expectation
Use models to simulate the road conditions of vehicles in the real world.
* # Video for simulation
[Video](https://drive.google.com/file/d/1MFCtU049tmFCPYqOp1yzHh61lf53hfTW/view?usp=share_link) 




\