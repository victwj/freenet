import matplotlib.pyplot as plt
import numpy as np

QuartileOne = []
QuartileThree = []
Median = []
# Mean = []

for i in range(10):

    fp = open("eval3/results/result{}.txt".format(i), "r")
    lines = [line.strip() for line in fp.readlines() if "code" in line or "Snapshot" in line]
    fp.close()

    lineCount = 0
    quartileOne = []
    quartileThree = []
    median = []
    # mean = []

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
                timestepLengths.append(500)
            lineCount += 1
        
        for i in range(100-len(timestepLengths)):
            timestepLengths.append(500)

        quartileOne.append(np.percentile(timestepLengths, 25))
        median.append(np.percentile(timestepLengths, 50))
        quartileThree.append(np.percentile(timestepLengths, 75))
        # mean.append(np.mean(timestepLengths))

        lineCount += 1

    QuartileOne.append(quartileOne)
    Median.append(median)
    QuartileThree.append(quartileThree)
    # Mean.append(mean)

QuartileOne = np.transpose(QuartileOne)
Median = np.transpose(Median)
QuartileThree = np.transpose(QuartileThree)
# Mean = np.transpose(Mean)

plot_y_q1 = [np.mean(x) for x in QuartileOne]
plot_y_m = [np.mean(x) for x in Median]
plot_y_q2 = [np.mean(x) for x in QuartileThree]
# plot_y_a = [np.mean(x) for x in Mean]
plot_x = [i for i in range(1, len(plot_y_m)+1)]

plt.plot(plot_x, plot_y_q1, linestyle=":", color="blue", label="first quartile")
plt.plot(plot_x, plot_y_m, linestyle="-", color="red", label="median")
plt.plot(plot_x, plot_y_q2, linestyle="-.", color="green", label="third quartile")
# plt.plot(plot_x, plot_y_a, linestyle="--", color="black", label="third quartile")

plt.yscale("log")
plt.xlabel("Node failure rate (%)")
plt.ylabel("Request pathlength (hops)")
plt.legend(loc=4)

plt.savefig("eval3/eval3.png")