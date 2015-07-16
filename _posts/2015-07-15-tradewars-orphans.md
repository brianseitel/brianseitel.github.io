---
layout: post
title: Remaking Tradewars - Part II
subtitle: Part II - Orphaned Sectors
date: 2015-07-16 11:23:45
disqus: n
---

### The Orphaned Sector Problem

In the [previous post](/2015/07/07/tradewars-big-bang/), we initiated the Big Bang and made some iterations to create a Universe containing Clusters of Sectors. At the end of the post, we ran into an issue whereupon we might wind up with orphaned sectors.

Here's a recap of the algorithm itself:

```
    # Define some constants
    define NUM_SECTORS = 1000
    define MIN_NEIGHBORS = 1
    define MAX_NEIGHBORS = 4
    define MIN_CLUSTERS  = 10
    define MIN_SECTORS_PER_CLUSTER = 2
    define MAX_SECTORS_PER_CLUSTER = 20

    # Begin Big Bang
    create empty Universe
    fill Universe with K Sectors
    create empty Clusters

    while still available Sectors
        add Cluster C to Clusters list
        for 0..K (where K = rng(MIN_SECTORS_PER_CLUSTER, MAX_CLUSTERS_PER_SECTOR)
            for 1..J (where J = rng(MIN_NEIGHBORS, MAX_NEIGHBORS)
                pick random available Sector as N (where N != S)
                add to S.Neighbors list
            add S to C.Sectors list
            remove S from available Sector list

        for each Cluster as C
        pick random Cluster N (where N != C)
            pick random Sector S1 in C
            pick random Sector S2 in N
            add S2 to S1.Neighbors list
```

For a simple Universe, you might wind up with Sectors in a Cluster like this:

<div id="cluster" style="width: 100%; height: 300px;"></div>

<div id="options"></div>

See the problem? Once we go from from S1 -> S2, we can't get back to S1 -- ever! Big problem!

How do we solve this? Well, that's the $64,000 question.

### Fagin's Revenge

The way to ensure that we have no orphaned sectors is to travel from each sector to every other sector. If we can make it from A -> B, A -> C, ..., A -> N, and then from B -> A, B -> C, ..., B -> N, and so on, then we have ensured that every Sector is reachable from every other Sector.

Right off the bat, though, you can tell that this is a *very* inefficient algorithm. For N sectors, we're going to check every other sector N times, for O(n^2) complexity.

There are two primary ways to search a graph like this. One is **depth-first search** and the other is **breadth-first search**. You can read the linked articles on Wikipedia for more information.

#### Depth-First Search

Here's a basic algorithm:

```
function dfs(graph, v):
  let S be a stack
  S.push(v)
  while S is not empty
    v = S.pop() 
    if v is not labeled as discovered:
      label v as discovered
      for all sectors from v to w in graph.adjacentEdges(v) do
        S.push(w)
```

In plain English, it works like this: Imagine you walk into a house. In the foyer, there are three neighboring rooms: the living room, the kitchen, and the office. You mentally mark the Foyer as "discovered." You pick one (the living room) and go there first. The living room has two neighboring rooms: the foyer and a closet. You mentally mark the living room as "discovered" and you enter the closet. There are no neighbors, so you back out to the living room. There are no undiscovered rooms, so you back up to the foyer. Now you've already visited the living room, so you move to the next neighboring room -- the kitchen. This process repeats itself.

By skipping rooms you've already discovered, you can prevent getting stuck into an infinite loop (foyer -> living room -> foyer -> living room -> foyer...).

The downside to this process is that the stack `S` can get very, very large. Depth-first search algorithms often suffer from non-termination -- that is, the universe can be so big that it just never finishes because it's almost impossible to traverse them all in good time.

Now let's look at Breadth-First Search.

#### Breadth-First Search

Here's a basic algorithm from Wikipedia:

```
function bfs(graph,v):
  let Q be a queue
  Q.enqueue(v)
  label v as discovered
  while Q is not empty
    v ← Q.dequeue()
    process(v)
    for all edges from v to w in G.adjacentEdges(v) do
      if w is not labeled as discovered
        Q.enqueue(w)
        label w as discovered
```

The first major difference is that we're using a queue instead of a stack. The difference is minor but important. A stack is like a stack of plates. In order to get to the bottom plate, you have to remove the top plates (pop them off). This is a last-in-first-out (LIFO) approach. A queue is more like a line at the bank. In order to get to the teller, you have to wait for all of the previously added people to finish. This is a first-in-first-out (FIFO) approach.

It turns out that the complexity of this is basically the same as DFS, but depending on the topography of the graph, the best case scenarios might be faster.

Here's a visual example of what it looks like: 

![Breadth-First Search](https://upload.wikimedia.org/wikipedia/commons/4/46/Animated_BFS.gif)

See how it evaluates each level before it goes on to the next? Contrast that with Depth-first search, which traverses all the way down the tree before going to the next branch.

#### Putting it all together

Now all we have to do is pick one of these algorithms, start the graph at Sector 1, conduct the search, and return a list of orphaned sectors. For each of these orphaned sectors, we connect them to another sector in the cluster, then continue on. It might look something like this:

```
  for C in Clusters
    for S in C.Sectors
      while (orphans = DFS(S))
        for O in orphans
          N = RNG(C.Sectors) // Pick a random sector in this cluster
          add O to N.Neighbors list
```

Be warned: this might take a LONG time for large numbers of sectors! However, once it's done, we should have a fully connected universe!

In the next post, we'll actually start writing some code.

<script type="text/javascript">
  // create an array with nodes
  var nodes = new vis.DataSet([
    {id: 1, label: 'Sector 1', color: '#CACACA'},
    {id: 2, label: 'Sector 2'},
    {id: 3, label: 'Sector 3'},
    {id: 4, label: 'Sector 4'},
    {id: 5, label: 'Sector 5'},
  ]);

  // create an array with edges
  var edges = new vis.DataSet([
    {from: 1, to: 2, arrows:'to'},
    {from: 1, to: 3, arrows:'to'},
    {from: 1, to: 4, arrows:'to'},
    {from: 2, to: 4, arrows:'to'},
    {from: 2, to: 5, arrows:'to'},
    {from: 3, to: 5, arrows:'to, from'},
    {from: 3, to: 2, arrows:'to, from'},
    {from: 4, to: 5, arrows:'to, from'},
  ]);

  // create a network
  var container = document.getElementById('cluster');
  var data = {
    nodes: nodes,
    edges: edges
  };
  var network = new vis.Network(container, data, options);
  var options = {
    interaction: { 
        dragNodes: false,
        dragView: false,
        zoomView: false,
        keyboard: false
    },
    layout: {
        randomSeed: 502998,
        hierarchical: {
            enabled: false,
            levelSeparation: 100,
            sortMethod: 'directed'
        }
    },
    physics: {
        enabled: true
    }
};
  network.setOptions(options);
</script>


