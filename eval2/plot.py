import matplotlib.pyplot as plt
import numpy as np

QuartileOne = []
QuartileThree = []
Median = []

for i in range(10):
    
    fp = open("eval2/results/result{}.txt".format(i), "r")
    lines = [line.strip() for line in fp.readlines() if "code" in line or "Snapshot" in line]
    fp.close()

    lineCount = 0
    quartileOne = []
    quartileThree = []
    median = []

    while lineCount < len(lines):

        while lineCount < len(lines) and "Start Snapshot" not in lines[lineCount]:
            lineCount += 1

        if lineCount == len(lines):
            break
        
        lineCount += 1
        timestepLengths = []

        while "End Snapshot" not in lines[lineCount]:
            msgCode = int(lines[lineCount].split(",")[0].split(":")[1])
            if msgCode == 21:
                timestepLengths.append(int(lines[lineCount].split(",")[1].split(":")[1]))
            else:
                timestepLengths.append(100)
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
plot_x = [20+(x*50) for x in range(1, len(plot_y_m)+1)]

plt.plot(plot_x, plot_y_q1, linestyle=":", color="blue", label="first quartile")
plt.plot(plot_x, plot_y_m, linestyle="-", color="red", label="median")
plt.plot(plot_x, plot_y_q2, linestyle="-.", color="green", label="third quartile")

plt.xscale("log")
plt.ylim(0, 120)
plt.xlabel("Network size (nodes)")
plt.ylabel("Request pathlength (hops)")
plt.legend(loc=2)

plt.savefig("eval2/eval2.png")