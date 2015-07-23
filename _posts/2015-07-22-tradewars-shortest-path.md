---
layout: post
title: Remaking Tradewars - Part IV
subtitle: Shortest Path Search
date: 2015-07-20 09:42:52
disqus: n
---

### Shortest Path

In the [previous post](/2015/07/20/tradewars-orphans-revisited/), we generated a fully-connected graph of the universe that looks something like this.

<div id="universe" style="width: 100%; height: 400px;"></div>

Great, we have a fully connected graph! 

In Tradewars, the goal is to travel from sector to sector, buying and selling goods. Some sectors might sell cheap Fuel, while others are willing to pay top dollar for Fuel. In that case, it would be wise to keep both sectors in mind for future reference. Next time you buy cheap fuel at Sector A, you can turn around and sell it at a high price at Sector B!

To make it slightly easier to travel from sector to sector, Tradewars offers a `jump` command, whereupon it spits out the shortest path between Sector A and Sector B, then automatically pushes you along that path, pausing in each sector to give you time to stop the jump if you see something interesting. Incidentally, this is also a great way to explore the universe!

The first obstacle in the `jump` command is finding the shortest path. There are a lot of ways to do this, but for now, we'll use something called [Dijkstra's Algorithm](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm). It's probably one of the better ways to find the shortest path in a graph like our universe, and it's relatively quick, too.

Here's the algorithm in pseudo-code, brazenly stolen from Wikipedia:

```
function Dijkstra(Graph, source):
    dist[source] ← 0  // Distance from source to source
    prev[source] ← undefined // Previous node in optimal path initialization

    for each vertex v in Graph:  // Initialization
       if v ≠ source:           // Where v has not yet been removed from Q (unvisited nodes)
           dist[v] ← infinity             // Unknown distance function from source to v
           prev[v] ← undefined            // Previous node in optimal path from source
       end if 
       add v to Q                     // All nodes initially in Q (unvisited nodes)
    end for

    while Q is not empty:
       u ← vertex in Q with min dist[u]  // Source node in first case
       remove u from Q 
       
       for each neighbor v of u:           // where v is still in Q.
           alt ← dist[u] + length(u, v)
           if alt < dist[v]:               // A shorter path to v has been found
               dist[v] ← alt 
               prev[v] ← u 
           end if
       end for
    end while
    return dist[], prev[]
end function
```

In graphical form, it looks something like this (also stolen from Wikipedia):

![Dijkstra's Algorithm](https://upload.wikimedia.org/wikipedia/commons/5/57/Dijkstra_Animation.gif)

There are some other variants of this, and I leave it as an experiment to the reader to try them out.

By passing in our Universe as the `Graph`, and the starting sector as the `source`, we can find all of the paths available. Finding the shortest path to a particular target is a matter of choosing the `target` item from the resulting `dist[]` that is returned from the function.

Tune in for Part V, when we talk about the next exciting module in Tradewars: outposts.

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