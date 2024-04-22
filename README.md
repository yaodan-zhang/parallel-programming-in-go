This is a thorough repository that showcases the parallel programming projects I did in Golang.

1. The first project is an image processing task, with a comparison to different program speeds using
   no parallelization / sequential implementation;
   parallelization among images;
   parallelization among image chunks.

2. The second project mimics a user's Twitter feed in managing his tweets
   a producer-consumer environment is used such that parallelization is realized in the number of consumers.

3. The third project reimplements the first project but uses a work-stealing confinement,
   that is, different threads have their own work queue filled with tasks by a producer, and if any of their queue is empty,
   the thread will try to steal task from another thread.

Enjoy reading the report analysis as it is somewhat insightful in figuring out how parallelization can be used to speed up
a program in some cases while not as expected in others.
