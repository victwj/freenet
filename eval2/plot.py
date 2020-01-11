import matplotlib.pyplot as plt
import numpy as np

QuartileOne = []
QuartileThree = []
Median = []

for i in range(10):

    fp = open("eval2/results/result{}.txt".format(i), "r")
    lines = [line.strip() for line in fp.readlines()]
    fp.close()

    lineCount = 0
    quartileOne = []
    quartileThree = []
    median = []

    while lineCount < len(lines) and "Node count 500 added" not in lines[lineCount]:
        lineCount += 1

    while lineCount < len(lines):

        currNodeCount = int(lines[lineCount].split(" ")[2])
        
        lineCount += 1
        timestepLengths = []

        while "Node count" not in lines[lineCount]:
            msgCode = int(lines[lineCount].split(",")[0].split(":")[1])
            if msgCode == 21:
                timestepLengths.append(int(lines[lineCount].split(",")[1].split(":")[1]))
            else:
                timestepLengths.append(500)
            lineCount += 1

        quartileOne.append(np.percentile(timestepLengths, 25))
        median.append(np.percentile(timestepLengths, 50))
        quartileThree.append(np.percentile(timestepLengths, 75))
        lineCount += 1

    QuartileOne.append(quartileOne)
    Median.append(median)
    QuartileThree.append(quartileThree)

QuartileOne = np.transpose(QuartileOne)
Median = np.transpose(Median)
QuartileThree = np.transpose(QuartileThree)

plot_y_q1 = [np.mean(x) for x in QuartileOne]
plot_y_m = [np.mean(x) for x in Median]
plot_y_q2 = [np.mean(x) for x in QuartileThree]

plt.plot(plot_y_q1, linestyle=":", color="blue", label="first quartile")
plt.plot(plot_y_m, linestyle="-", color="red", label="median")
plt.plot(plot_y_q2, linestyle="-.", color="green", label="third quartile")

plt.xscale("log")
plt.xlim(100, 1000000)
plt.xlabel("Network size (nodes)")
plt.ylabel("Request pathlength (hops)")
plt.legend()

plt.savefig("eval2/eval2.png")