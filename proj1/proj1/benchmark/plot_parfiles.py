import matplotlib.pyplot as plt

for data in ["small", "mixture", "big"]:
    threads = [2,4,6,8,12]
    speedup = []
    with open(data + "_parfiles_1.txt") as f:
        benchmark = float(f.read().strip())
    for t in threads:
        with open(data + "_parfiles_" + str(t) + ".txt") as f:
            speedup.append (benchmark/float(f.read().strip()))
    plt.plot(threads, speedup, label = data)
       
plt.xlabel('Number of Threads')
plt.ylabel('Speedup')
plt.legend()
plt.savefig('speedup-images.png')