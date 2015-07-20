---
layout: post
title: Remaking Tradewars - Part III
subtitle: Ensuring Connectivity
date: 2015-07-20 09:42:52
disqus: n
---

### The Orphaned Sector Problem Revisited

In the [first post](/2015/07/07/tradewars-big-bang/), we initiated the Big Bang and made some iterations to create a Universe containing Clusters of Sectors. At the end of the post, we ran into an issue whereupon we might wind up with orphaned sectors.

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

<div id="orphans" style="width: 100%; height: 300px;"></div>

<div id="options"></div>

See the problem? Once we go from from S1 -> S2, we can't get back to S1 -- ever! Big problem!

In [Part II](/2015-07/15-tradewars-orphan), I discussed a possible solution: using depth-first search (DFS) or breadth-first search (BFS) to detect orphan sectors, then attempt to connect them again.

It turns out, there's a better way.

### Tunnel Vision

This one took me awhile to figure out. As I was tinkering with the algorithms to get rid of the orphan sectors, I realized something: I was randomly connecting previously-randomized sectors. Why bother with the random connections if the sector numbers are already random? Just connect it to the next one in the list!

In other words, we have some clusters that look like this:


```
sectors = [1,2,3,4,5,6,7,8,9,10];
cluster1 = [6,3,1,7];
cluster2 = [2,9,4];
cluster3 = [10,5,8];
```

Since each cluster is already randomized, let's just connect them in order!

```
cluster1: 6 -> 3 -> 1 -> 7
cluster2: 2 -> 9 -> 4
cluster3: 10 -> 5 -> 8
```

Once each cluster's sectors are connected, connect the first and last sector of each cluster:
```
cluster1: 8 (cluster3) -> 6 -> 3 -> 1 -> 7 -> 2 (cluster2)
cluster2: 7 (cluster1) -> 2 -> 9 -> 4 -> 10 (cluster3)
cluster3: 4 (cluster2) -> 10 -> 5 -> 8 -> 6 (cluster1)
```

If you were to graph this, you'd get a circle, like this:

<div id="circle" style="width: 100%; height: 400px;"></div>

As you can see, it's possible to get from any sector to any other sector in the universe. Now that we know that's possible, all we have to do is *add new connections* at random within clusters. The algorithm looks something like this:


```
// First pass: Connect each sector sequentially
for C in Clusters
  for S in C.Sectors
    Neighbor = S + 1
    if (Neighbor > number of sectors in C)
      Neighbor = 0
    endif

    S.connect(Neighbor)
  endfor

  NeighborC = C + 1
  if (NeighborC > number of clusters)
    NeighborC = 0
  endif

  // Connect last sector of this cluster to 
  // first sector of next cluster
  C.Sectors[last].connect(NeighborC.Sectors[0])
endfor

// Second pass: Randomly connect sectors
for C in Clusters
  for S1 in C.Sectors
    chance = Random(0, 100)
    if (Random > CHANCE_OF_CONNECTION)
      S2 = C.Sectors.Random
      S1.connect(S2)
    endif
  endfor

  // Select random cluster
  C2 = Clusters.Random

  // Connect random sector in C to random sector in C2
  C.connect(C2);
endfor

```

This might look something like this:

<div id="universe" style="width: 100%; height: 400px;"></div>

If you look at this carefully, you'll see that there are no orphaned nodes. The universe is 100% connected, but still looks pretty random!

If you're interested in exploring this a little more deeply, try modifying your algorithm to implement the following ideas:

  * Add more sectors!
  * Use [greuler](maurizzzio.github.io/greuler/) or similar tool to generate a visual graph of your universe, like my examples above
  * Introduce dead-ends (e.g., there's only one path to a node)
    * Example: 2 <-> 6 <-> 8 (no other path to 8 but through 2 and 6)
  * Introduce tunnels (e.g., only one path between A and B, but many paths to A and many paths from B)
    * Example: 1,2,3 -> 4 -> 5 -> 6 -> 7,8,9 (no other path to 7 except through 4, 5, and 6)
  * Write a test using DFS or BFS to ensure all nodes are reachable from every other node
  * Write a function using Dijkstra's Algorithm to determine the shortest path from A to B, where A and B are different sectors. (*Note: We will discuss this in a future post!*)

Tune in next time for more Tradewars-esque fun!

Toodles!


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
  var container = document.getElementById('orphans');
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

<script type="text/javascript">
  // create an array with nodes
  var nodes = new vis.DataSet([
    {id: 1, label: 'Sector 1'},
    {id: 2, label: 'Sector 2'},
    {id: 3, label: 'Sector 3'},
    {id: 4, label: 'Sector 4'},
    {id: 5, label: 'Sector 5'},
    {id: 6, label: 'Sector 6'},
    {id: 7, label: 'Sector 7'},
    {id: 8, label: 'Sector 8'},
    {id: 9, label: 'Sector 9'},
    {id: 10, label: 'Sector 10'},

  ]);

  // create an array with edges
  var edges = new vis.DataSet([
    {from: 1, to: 7, arrows:'to'},
    {from: 2, to: 9, arrows:'to'},
    {from: 3, to: 1, arrows:'to'},
    {from: 4, to: 10, arrows:'to'},
    {from: 5, to: 8, arrows:'to'},
    {from: 6, to: 3, arrows:'to'},
    {from: 7, to: 2, arrows:'to'},
    {from: 8, to: 6, arrows:'to'},
    {from: 9, to: 4, arrows:'to'},
    {from: 10, to: 5, arrows:'to'},
  ]);

  // create a network
  var container = document.getElementById('circle');
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

<script type="text/javascript">
  // create an array with nodes
  var nodes = new vis.DataSet([
    {id: 1, label: 'Sector 1'},
    {id: 2, label: 'Sector 2'},
    {id: 3, label: 'Sector 3'},
    {id: 4, label: 'Sector 4'},
    {id: 5, label: 'Sector 5'},
    {id: 6, label: 'Sector 6'},
    {id: 7, label: 'Sector 7'},
    {id: 8, label: 'Sector 8'},
    {id: 9, label: 'Sector 9'},
    {id: 10, label: 'Sector 10'},

  ]);

/*
cluster1: 8 (cluster3) -> 6 -> 3 -> 1 -> 7 -> 2 (cluster2)
cluster2: 7 (cluster1) -> 2 -> 9 -> 4 -> 10 (cluster3)
cluster3: 4 (cluster2) -> 10 -> 5 -> 8 -> 6 (cluster1)
 */
  // create an array with edges
  var edges = new vis.DataSet([
    {from: 1, to: 7, arrows:'to'},
    {from: 1, to: 3, arrows:'to'},
    {from: 1, to: 2, arrows:'to'},
    {from: 2, to: 9, arrows:'to'},
    {from: 2, to: 6, arrows:'to'},
    {from: 3, to: 1, arrows:'to'},
    {from: 4, to: 10, arrows:'to'},
    {from: 4, to: 9, arrows:'to'},
    {from: 4, to: 2, arrows:'to'},
    {from: 5, to: 8, arrows:'to'},
    {from: 5, to: 6, arrows:'to'},
    {from: 6, to: 3, arrows:'to'},
    {from: 6, to: 2, arrows:'to'},
    {from: 6, to: 1, arrows:'to'},
    {from: 7, to: 2, arrows:'to'},
    {from: 7, to: 6, arrows:'to'},
    {from: 8, to: 6, arrows:'to'},
    {from: 8, to: 5, arrows:'to'},
    {from: 8, to: 4, arrows:'to'},
    {from: 9, to: 4, arrows:'to'},
    {from: 9, to: 2, arrows:'to'},
    {from: 10, to: 5, arrows:'to'},
  ]);

  // create a network
  var container2 = document.getElementById('universe');
  var data2 = {
    nodes: nodes,
    edges: edges
  };
  var network = new vis.Network(container2, data2, options);
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


