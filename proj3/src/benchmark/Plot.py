import subprocess
import matplotlib.pyplot as plt
threads = [2,4,6,8,12]

# Speedup graph for parfiles
def plot_parfiles():
    speedup = []
    # Run sequential ver
    parfileS = 0
    subprocess.run(["touch","parfileS.txt"])
    with open("parfileS.txt","a") as outfile:
        subprocess.run(["go", "run","../editor/editor.go","test","parfiles","1"], stdout = outfile)
    with open("parfileS.txt","r") as infile:
        parfileS = float(infile.readline().strip())
    subprocess.run(["rm","parfileS.txt"])
    # Run parallel ver
    for t in threads:
        subprocess.run(["touch","parfile"+str(t)+"T"+".txt"])
        with open("parfile"+str(t)+"T"+".txt","a") as outfile:
            subprocess.run(["go","run","../editor/editor.go","test","parfiles",str(t)], stdout = outfile)
        parfileT = 0
        with open("parfile"+str(t)+"T"+".txt","r") as infile:
            parfileT = float(infile.readline().strip())
        speedup.append(parfileS/parfileT)
        subprocess.run(["rm","parfile"+str(t)+"T"+".txt"])

    # Plot 
    plt.plot(threads, speedup, label = "parfiles")
    
# Speedup graph for parslices
def plot_parslices():
    speedup = []
    # Run sequential ver
    subprocess.run(["touch","parsliceS.txt"])
    with open("parsliceS.txt","a") as outfile:
        subprocess.run(["go", "run","../editor/editor.go","test","parslices","1"], stdout = outfile)
    parsliceS = 0
    with open("parsliceS.txt","r") as infile:
        parsliceS = float(infile.readline().strip())
    subprocess.run(["rm","parsliceS.txt"])

    # Run parallel ver
    for t in threads:
        subprocess.run(["touch","parslice"+str(t)+"T"+".txt"])
        with open("parslice"+str(t)+"T"+".txt","a") as outfile:
            subprocess.run(["go","run","../editor/editor.go","test","parslices",str(t)], stdout = outfile)
        parsliceT = 0
        with open("parslice"+str(t)+"T"+".txt","r") as infile:
            parsliceT = float(infile.readline().strip())
        speedup.append(parsliceS/parsliceT)
        subprocess.run(["rm","parslice"+str(t)+"T"+".txt"])

    # Plot 
    plt.plot(threads, speedup, label = "parslices")

def main():
    plot_parfiles()
    plot_parslices()
    #Plot
    plt.xlabel('Number of Threads')
    plt.ylabel('Speedup')
    plt.legend()
    plt.savefig('speedup-image.png')

main()
