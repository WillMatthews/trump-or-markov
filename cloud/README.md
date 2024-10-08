# Real & Fake Trump Tweet Generator API

Currently all Markov Chains & Tweets are stored in-memory which is not great.
Training requires significant CPU but it's only done once.
Inference costs nothing (other than holding maps in memory, which I admit is extremely expensive).

I did a crude memory usage measurement. I think I can do better than this.

| Object       | Memory in MB | Memory in MB (hash key) |
|--------------|--------------|-------------------------|
| Binary       | 9.9          | 9.9                     |
| Tweets Alone | 33.1         | 33.1                    |
| Ord 1 chain  | 2            | 58                      |
| Ord 2 chain  | 180          | 80                      |
| Ord 3 chain  | 110          | 93                      |
| Ord 4 chain  | 255          | 171                     |

I probably shouldn't hash the strings before storing in map. The map should just do that for me efficiently.

Potential strategies for reducing memory usage:
- Use pointers liberally and point at a dictionary of words
- Use a trie instead of a map (? I think this might work / be more memory efficient)
- Store frequency counts for each word rather than repeating the word in the `possibles` slice.
