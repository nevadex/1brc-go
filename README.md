# 7.3 seconds

## this haunted my brain while trying to sleep for  i think 4 nights.

i have a lot to learn

---

### key takeaways:

- [locks are bad](https://x.com/TetraspaceWest/status/1800354963337441646) (too expensive for high perf, low latency apps, i was seeing minutes spent waiting for a lock to open for writing)
- maps are slower than i thought, i think i confused them non-integer arrays (which i wish upon a star for).
- finally understood how go concurrency works. theres a GT demo somewhere that failed due to race conditions in websocket handlers :(
- pprof is amazing, i love it sm (helped me understand a lot about what was slowing the program down, but the actual numbers were much exaggerated, dunno why still)
- binary manipulation is actually easier than i thought, learnt how bitwise ops work (dark magic in our 204 rpi pico CAN code made by an electrical engineer)
- pre-alloc works wonders
- sync.wg and sync.pool are great (pool isnt for this though, floats are smol)
- for loops are kind of slow.. (ASM indexbyte func was faster, surprised me for some reason)
- floats arent too expensive
- speaking of floats, wasted HOURS learning about how floats are mapped in memory, only to realize that i DO NOT need binary floats, i could just do the stupid per-digit thing
- speaking of floats again, emailed librarian to get access to IEEE Xplore (should be free but wtv)
- go testing framework is FIRE, benchmarks are really simple and work well
- array append is cheaper than i thought (15ns vs 21ns)

## and worst of all,

comparison is the true destroyer of joy (do not look at 1brc java lbs :(

### over the entire project, i think the time went from estimated 20 hrs, to 10 minutes, to 7.3 seconds. proud of it 

---

## how to run:

theres like 20 goland run configs and stuff in the .idea folder (different tests, file sizes)

could also just `go run .` just make sure the path in the main is correct