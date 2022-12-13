from ctypes import cdll
import ctypes



class Car(ctypes.Structure):
    _fields_ = [
        ("speed", ctypes.c_int),
        ("kind", ctypes.c_int),
        ("backlight", ctypes.c_int),
        ("turninglight", ctypes.c_int),
        ("accel", ctypes.c_int),
    ]




PlayMultiLaneModel = cdll.LoadLibrary('./s1.so').PlayMultiLaneModel
initialMultiRoad = cdll.LoadLibrary('./s1.so').initialMultiRoad
return_int_array = cdll.LoadLibrary('./s1.so').returnIntArray
args_type = ctypes.POINTER(ctypes.c_int64) * 2


numGens = 1000

cellWidth = 21

trafficLightLane = [1,3]
trafficLightPos = 25
trafficLightTime = [0,0,0]
trafficLightTime[0] = 30
trafficLightTime[1] = 5  
trafficLightTime[2] = 30 

sdvPercentage = 0.5
nsdvPercentage = 1.0 - sdvPercentage
laneNum = 5





initialMultiRoad.argtypes= [args_type,ctypes.c_int,ctypes.c_int]
initialmultiRoad = initialMultiRoad(trafficLightLane, trafficLightPos, laneNum)
timePointsMulti, totalCnt = PlayMultiLaneModel(initialMultiRoad, numGens, trafficLightPos, laneNum, trafficLightLane, trafficLightTime, nsdvPercentage)

