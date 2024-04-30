#!/bin/bash
#
#SBATCH --mail-user=katherinezh@cs.uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj1_benchmark 
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/katherinezh/project-1-yaodan-zhang/proj1/benchmark/
#SBATCH --partition=debug 
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=4:00:00


module load golang/1.19
# Test parfiles.go:
# Run small data set
go run ../editor/editor.go small parfiles 1 > small_parfiles_1.txt
go run ../editor/editor.go small parfiles 2 > small_parfiles_2.txt
go run ../editor/editor.go small parfiles 4 > small_parfiles_4.txt
go run ../editor/editor.go small parfiles 6 > small_parfiles_6.txt
go run ../editor/editor.go small parfiles 8 > small_parfiles_8.txt
go run ../editor/editor.go small parfiles 12 > small_parfiles_12.txt
# Run mixture data set
go run ../editor/editor.go mixture parfiles 1 > mixture_parfiles_1.txt
go run ../editor/editor.go mixture parfiles 2 > mixture_parfiles_2.txt
go run ../editor/editor.go mixture parfiles 4 > mixture_parfiles_4.txt
go run ../editor/editor.go mixture parfiles 6 > mixture_parfiles_6.txt
go run ../editor/editor.go mixture parfiles 8 > mixture_parfiles_8.txt
go run ../editor/editor.go mixture parfiles 12 > mixture_parfiles_12.txt
# Run big data set
go run ../editor/editor.go big parfiles 1 > big_parfiles_1.txt
go run ../editor/editor.go big parfiles 2 > big_parfiles_2.txt
go run ../editor/editor.go big parfiles 4 > big_parfiles_4.txt
go run ../editor/editor.go big parfiles 6 > big_parfiles_6.txt
go run ../editor/editor.go big parfiles 8 > big_parfiles_8.txt
go run ../editor/editor.go big parfiles 12 > big_parfiles_12.txt
# Generate speedup plots
python plot_parfiles.py
# Delete all intermediate files
rm -r *.txt
# Delete all pictures
rm ../data/out/*
# Test parslices.go:
# Run small data set
go run ../editor/editor.go small parslices 1 > small_parslices_1.txt
go run ../editor/editor.go small parslices 2 > small_parslices_2.txt
go run ../editor/editor.go small parslices 4 > small_parslices_4.txt
go run ../editor/editor.go small parslices 6 > small_parslices_6.txt
go run ../editor/editor.go small parslices 8 > small_parslices_8.txt
go run ../editor/editor.go small parslices 12 > small_parslices_12.txt
# Run mixture data set
go run ../editor/editor.go mixture parslices 1 > mixture_parslices_1.txt
go run ../editor/editor.go mixture parslices 2 > mixture_parslices_2.txt
go run ../editor/editor.go mixture parslices 4 > mixture_parslices_4.txt
go run ../editor/editor.go mixture parslices 6 > mixture_parslices_6.txt
go run ../editor/editor.go mixture parslices 8 > mixture_parslices_8.txt
go run ../editor/editor.go mixture parslices 12 > mixture_parslices_12.txt
# Run big data set
go run ../editor/editor.go big parslices 1 > big_parslices_1.txt
go run ../editor/editor.go big parslices 2 > big_parslices_2.txt
go run ../editor/editor.go big parslices 4 > big_parslices_4.txt
go run ../editor/editor.go big parslices 6 > big_parslices_6.txt
go run ../editor/editor.go big parslices 8 > big_parslices_8.txt
go run ../editor/editor.go big parslices 12 > big_parslices_12.txt
# Generate speedup plots
python plot_parslices.py
# Delete all intermediate files
rm -r *.txt
