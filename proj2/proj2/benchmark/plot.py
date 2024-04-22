import subprocess
import matplotlib.pyplot as plt
threads = [2,4,6,8,12]

# Speedup graph for xsmall
def plot_xsmall():
    speedup = []
    # Run sequential ver
    subprocess.run(["touch","xsmallS.txt"])
    with open("xsmallS.txt","a") as outfile:
        for i in range(5):
            subprocess.run(["go", "run","benchmark.go","s","xsmall"], stdout = outfile)
    xsmallS = 0
    with open("xsmallS.txt","r") as infile:
        for i in range(5):
            xsmallS += float(infile.readline().strip())
    xsmallS /= 5
    subprocess.run(["rm","xsmallS.txt"])

    # Run parallel ver
    for t in threads:
        subprocess.run(["touch","xsmall"+str(t)+"T"+".txt"])
        with open("xsmall"+str(t)+"T"+".txt","a") as outfile:
            for i in range(5):
                subprocess.run(["go","run","benchmark.go","p","xsmall",str(t)], stdout = outfile)
        xsmallT = 0
        with open("xsmall"+str(t)+"T"+".txt","r") as infile:
            for i in range(5):
                xsmallT += float(infile.readline().strip())
        xsmallT /= 5
        speedup.append(xsmallS/xsmallT)
        subprocess.run(["rm","xsmall"+str(t)+"T"+".txt"])

    # Plot 
    plt.plot(threads, speedup, label = "xsmall")
    
# Speedup graph for small
def plot_small():
    speedup = []
    # Run sequential ver
    subprocess.run(["touch","smallS.txt"])
    with open("smallS.txt","a") as outfile:
        for i in range(5):
            subprocess.run(["go", "run","benchmark.go","s","small"], stdout = outfile)
    smallS = 0
    with open("smallS.txt","r") as infile:
        for i in range(5):
            smallS += float(infile.readline().strip())
    smallS /= 5
    subprocess.run(["rm","smallS.txt"])

    # Run parallel ver
    for t in threads:
        subprocess.run(["touch","small"+str(t)+"T"+".txt"])
        with open("small"+str(t)+"T"+".txt","a") as outfile:
            for i in range(5):
                subprocess.run(["go","run","benchmark.go","p","small",str(t)], stdout = outfile)
        smallT = 0
        with open("small"+str(t)+"T"+".txt","r") as infile:
            for i in range(5):
                smallT += float(infile.readline().strip())
        smallT /= 5
        speedup.append(smallS/smallT)
        subprocess.run(["rm","small"+str(t)+"T"+".txt"])

    # Plot 
    plt.plot(threads, speedup, label = "small")

# Speedup graph for medium
def plot_medium():
    speedup = []
    # Run sequential ver
    subprocess.run(["touch","mediumS.txt"])
    with open("mediumS.txt","a") as outfile:
        for i in range(5):
            subprocess.run(["go", "run","benchmark.go","s","medium"], stdout = outfile)
    mediumS = 0
    with open("mediumS.txt","r") as infile:
        for i in range(5):
            mediumS += float(infile.readline().strip())
    mediumS /= 5
    subprocess.run(["rm","mediumS.txt"])

    # Run parallel ver
    for t in threads:
        subprocess.run(["touch","medium"+str(t)+"T"+".txt"])
        with open("medium"+str(t)+"T"+".txt","a") as outfile:
            for i in range(5):
                subprocess.run(["go","run","benchmark.go","p","medium",str(t)], stdout = outfile)
        mediumT = 0
        with open("medium"+str(t)+"T"+".txt","r") as infile:
            for i in range(5):
                mediumT += float(infile.readline().strip())
        mediumT /= 5
        speedup.append(mediumS/mediumT)
        subprocess.run(["rm","medium"+str(t)+"T"+".txt"])

    # Plot 
    plt.plot(threads, speedup, label = "medium")

# Speedup graph for large
def plot_large():
    speedup = []
    # Run sequential ver
    subprocess.run(["touch","largeS.txt"])
    with open("largeS.txt","a") as outfile:
        for i in range(5):
            subprocess.run(["go", "run","benchmark.go","s","large"], stdout = outfile)
    largeS = 0
    with open("largeS.txt","r") as infile:
        for i in range(5):
            largeS += float(infile.readline().strip())
    largeS /= 5
    subprocess.run(["rm","largeS.txt"])

    # Run parallel ver
    for t in threads:
        subprocess.run(["touch","large"+str(t)+"T"+".txt"])
        with open("large"+str(t)+"T"+".txt","a") as outfile:
            for i in range(5):
                subprocess.run(["go","run","benchmark.go","p","large",str(t)], stdout = outfile)
        largeT = 0
        with open("large"+str(t)+"T"+".txt","r") as infile:
            for i in range(5):
                largeT += float(infile.readline().strip())
        largeT /= 5
        speedup.append(largeS/largeT)
        subprocess.run(["rm","large"+str(t)+"T"+".txt"])

    # Plot 
    plt.plot(threads, speedup, label = "large")

# Speedup graph for xlarge
def plot_xlarge():
    speedup = []
    # Run sequential ver
    subprocess.run(["touch","xlargeS.txt"])
    with open("xlargeS.txt","a") as outfile:
        for i in range(5):
            subprocess.run(["go", "run","benchmark.go","s","xlarge"], stdout = outfile)
    xlargeS = 0
    with open("xlargeS.txt","r") as infile:
        for i in range(5):
            xlargeS += float(infile.readline().strip())
    xlargeS /= 5
    subprocess.run(["rm","xlargeS.txt"])

    # Run parallel ver
    for t in threads:
        subprocess.run(["touch","xlarge"+str(t)+"T"+".txt"])
        with open("xlarge"+str(t)+"T"+".txt","a") as outfile:
            for i in range(5):
                subprocess.run(["go","run","benchmark.go","p","xlarge",str(t)], stdout = outfile)
        xlargeT = 0
        with open("xlarge"+str(t)+"T"+".txt","r") as infile:
            for i in range(5):
                xlargeT += float(infile.readline().strip())
        xlargeT /= 5
        speedup.append(xlargeS/xlargeT)
        subprocess.run(["rm","xlarge"+str(t)+"T"+".txt"])

    # Plot 
    plt.plot(threads, speedup, label = "xlarge")

def main():
    # plot speedups for all sizes
    plot_xsmall()
    plot_small()
    plot_medium()
    plot_large()
    plot_xlarge()

    # Label the axis of the plot and save it as an image
    plt.xlabel('Number of Threads')
    plt.ylabel('Speedup')
    plt.legend()
    plt.savefig('speedup-image.png')

main()
