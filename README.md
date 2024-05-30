Mosaic Collage Generator
Muhammad Aulia Firmansyah
mauliafirmansyah@uchicago.edu

USAGE
  -B float
        intensity of tile images color blend-in in float (0.0 - 1.0). 1.0 for full blend in, 0.0 for tiles image only (default 0.8)
  -I float
        intensity of mosaic images in float (0.0 - 1.0). 1.0 for full mosaic images, 0.0 for input image only (default 0.8)
  -M string
        running mode: s=sequential(default), p=parallel, w=parallel with work steal (default "s")
  -T int
        Number of goroutines. ignored if sequential. Must be positive (default 1)
  -U int
        Input image upscaling in integer. Must be positive (default 1)
  -d string
        Path to the mosaic tiles directory
  -i string
        Path to the input image
  -o string
        Path to the output image
  -s int
        Size of mosaic tiles in pixels. Must be positive

HOW TO RUN
To run this project, go to “proj3” folder, then run “python3 benchmark/benchmark-proj3.py”. To ensure that the script run without problem, Python version 3 should be installed, with matplotlib package included. Before running the script, the dataset should be put into its respective folder. The dataset can be accessed using this link: proj3-muhauliaf-extra

RESULT
After running the script file, there should be several outputs produced, which are the followings:
- Image output, which are the generated images from mosaic collage generator. There is one image output for each test size (small, medium, large)
- The benchmark-proj3.json file, which contains runtimes for both parts, all versions (sequential, parallel, work stealing), all sizes (small, medium, large), and all threads for parallel versions.
- Graph images, which shows speedups of both parallel versions (parallel, work stealing) compared to sequential version for both parts.