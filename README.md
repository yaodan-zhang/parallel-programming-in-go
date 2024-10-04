This repository showcases the parallel programming projects I did in Go.

1. The first project is an image processing task, with a comparison among different program speeds using
   no parallelization / sequential implementation;
   parallelization across images;
   parallelization across image chunks.

2. The second project mimics a user's Twitter feed in managing the tweets, and a producer-consumer environment is used such that parallelization is realized in the choice of the number of consumers.

3. The third project reimplements the first project but uses a work-stealing confinement, that is, different threads have their own work queue getting filled with tasks by a producer, and if any of the queues is empty, the thread will try to randomly steal tasks from another thread.

Enjoy reading the report analysis as it is somewhat insightful in figuring out how parallelization can be used to speed up a program in some cases while not as expected in others.

All copyrights are reserved to the Parallel Programming course from University of Chicago, which is where the projects are originated, and the author, who is the implementer.
