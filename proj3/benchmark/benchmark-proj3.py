import json
from pprint import pprint
import subprocess
import os
import matplotlib.pyplot as plt

# Predefined configurations
benchmark_output_file = "benchmark/benchmark-proj3.json"
benchmark_graph_dir = "benchmark/graph"
input_file = "data/in/phoenix.png"
output_file = "data/out/phoenix"
parallel_threads = [2,4,6,8,12]
mosaic_intensity = 0.7
color_blending = 0.7
test_times = 5
test_sizes = ["small","medium","large"]

# Configuration of various test suites
test_suites = {
    "small":{
        "tiles_dir":"data/tiles/small",
        "tile_size":50,
        "upscale":8,
    }, "medium":{
        "tiles_dir":"data/tiles/medium",
        "tile_size":200,
        "upscale":16,
    }, "large":{
        "tiles_dir":"data/tiles/large",
        "tile_size":400,
        "upscale":32,
    },
}

def main():
    if os.path.exists(benchmark_output_file):
        with open(benchmark_output_file, "r") as json_file:
            results = json.load(json_file)
    else:
        results = run_complete_benchmark()
        with open(benchmark_output_file, "w") as json_file:
            json.dump(results, json_file)
    print("Results data:")
    pprint(results)
    plot_results(results)

# Runs a complete benchmark for all test cases
def run_complete_benchmark():
    results = {}
    print("Running sequential version..")
    result_sequential = {}
    for test_size in test_sizes:
        result = []
        for i in range(test_times):
            result.append(run_benchmark(test_size,"s"))
        result_sequential[test_size] = result
    results["sequential"] = result_sequential
    print("Running parallel version..")
    result_parallel = {}
    for test_size in test_sizes:
        result_test = {}
        for thread in parallel_threads:
            result = []
            for i in range(test_times):
                result.append(run_benchmark(test_size,"p",thread))
            result_test[thread] = result
        result_parallel[test_size] = result_test
    results["parallel"] = result_parallel
    print("Running work stealing version..")
    result_parallel = {}
    for test_size in test_sizes:
        result_test = {}
        for thread in parallel_threads:
            result = []
            for i in range(test_times):
                result.append(run_benchmark(test_size,"w",thread))
            result_test[thread] = result
        result_parallel[test_size] = result_test
    results["worksteal"] = result_parallel
    return results

# Runs benchmark on a single configuration
def run_benchmark(test_size, run_mode, thread_count=1):
    test_suite = test_suites[test_size]
    tiles_directory = test_suite["tiles_dir"]
    tile_size = test_suite["tile_size"]
    upscale = test_suite["upscale"]
    command =  (f"go run proj3/mosaic "
                f"-i {input_file} "
                f"-o {output_file}-{test_size}.png "
                f"-d {tiles_directory} "
                f"-s {tile_size} "
                f"-U {upscale} "
                f"-I {mosaic_intensity} "
                f"-B {color_blending} "
                f"-M {run_mode} "
                f"-T {thread_count} "
    )
    print(f"Running: {command}")
    gorun = subprocess.run(command, shell=True, capture_output=True, text=True)
    run_result = [float(x) for x in gorun.stdout.strip().split()]
    print(f"Result: {run_result}")
    return run_result

# Plots results to several charts
def plot_results(results):
    parts = list(range(2))
    versions = ["parallel","worksteal"]
    speedups = {}
    for part in parts:
        speedups_part = {}
        for version in versions:
            speedups_version = {}
            for test_size in test_sizes:
                speedup_test = {}
                result_values = [x[part] for x in results["sequential"][test_size]]
                result_avg = sum(result_values)/len(result_values) if len(result_values)>0 else 0
                sequential_time = result_avg
                for thread in parallel_threads:
                    results[version][test_size] = {int(k):v for k,v in results[version][test_size].items()}
                    result_values = [x[part] for x in results[version][test_size][thread]]
                    result_avg = sum(result_values)/len(result_values) if len(result_values)>0 else 0
                    parallel_time = result_avg
                    speedup_test[thread] = sequential_time/parallel_time
                speedups_version[test_size] = speedup_test
            speedups_part[version] = speedups_version
        speedups[part] = speedups_part
    print("Speedups data:")
    pprint(speedups)

    for part in parts:
        for version in versions:
            speedups_part_version = speedups[part][version]
            x_values = parallel_threads
            y_values = {test: [speedups_part_version[test][thread] for thread in parallel_threads] for test in test_sizes}

            plt.figure(figsize=(8, 6))
            for legend, y_data in y_values.items():
                plt.plot(x_values, y_data, marker='o', label=legend)

            # Adding title and labels
            plt.title(f"Speedup Graph for part {part+1} of {version} version")
            plt.xlabel("Number of Threads")
            plt.ylabel("Speedup")

            # Adding legends
            plt.legend(title="Lines")

            plt.xticks(x_values)

            # Set y-axis ticks to show only 0.01 decimal
            plt.gca().yaxis.set_major_formatter(plt.FormatStrFormatter('%.2f'))

            # Save the graph to a file
            graph_path = f"{benchmark_graph_dir}/proj3-{version}-part{part+1}-graph.png"
            plt.savefig(graph_path)

            # Display the file path
            print("Graph saved to:", graph_path)

if __name__ == "__main__":
    main()