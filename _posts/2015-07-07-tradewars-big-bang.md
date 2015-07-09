---
layout: post
title: Remaking Tradewars - The Big Bang
date: 2015-07-07 23:01:33
disqus: n
---

### The Story

Back in *my* day, before new-fangled AOL and Prodigy, before broadband and 56kbps modems, there was the wonderful, wild-west world of the BBS (Bulletin Board System). Anyone with a computer and a phone line could set one up. Others could dial in, post on message forums and play games. Oh, games. The games.

L.O.R.D., Barren Realms Elite, and Tradewars. Tradewars was by far my favorite.

The premise was rather simple: you're a trader in a spaceship. The goal was to traverse the universe, buying and selling goods to accumulate money, or credits. With more credits, you could upgrade your ship, buy more holds (to carry more goods) or more weapons (to blow up other traders or aliens), and even buy your own planet and build a citadel!

It is with great nostalgia that I think back upon that game. It was text-only, with a sprinkle of ASCII color, but it captured my imagination.

Recently, I got to thinking about that game again. I had also gotten a hankering to learn [Go](http://golang.org) and this seemed like a great opportunity.

### The Big Bang

The first thing that happens when you set up a new Tradewars universe is The Big Bang. It's a quick program that instantiates the entire universe -- every sector, every cluster, every trading post, alien outpost, and planet.

You can think of each sector as a linked list. It contains its own data, plus the neighboring sectors. It might look something like this:


#### Sector 100

<table class="table table-condensed" style="border: 1px solid #CCC;width: 40%;margin: 0 60px;">
    <thead>
        <th>Name</th>
        <th>Neighbors</th>
    </thead>
    <tbody>
        <td>
            Terra IV
        </td>
        <td>
            <ul>
                <li>153</li>
                <li>6699</li>
                <li>981</li>
            </ul>
        </td>
    </tbody>
</table>

Each sector has a **Name**, **Number**, and list of **Neighbors**. Once we load a particular Sector, it's trivial to get to the next -- just select a neighbor and load it up!

Now, the hard part is generating this. Let's try a naive algorithm and see what happens.

### First Pass Big Bang

```
    create empty Universe
    fill Universe with W Sectors
    for each Sector as S:
        for between 1 and 4 (random number of neighbors)
            pick random Sector as N (where N != S)
            add to S.Neighbors list
```

OK! So with this basic algorithm we've populated the entire universe, randomly connecting sectors with each other! Great!

Except... (you knew there was a "but", right?) This presents a handful of problems.

### The Problems


#### There Goes the Neighborhood
First, let's imagine the real universe. In our solar system, Earth is between Mars and Venus. On the other hand, Uranus is over 2.5 billion kilometers away from Earth. That's not even close. Would it make sense to say Venus and Earth are neighbors? Maybe they belong in the same community (we'll get to this later), but they are certainly not neighbors.

With our naive algorithm above, we're randomly selecting from W sectors. Imagine a 1000 Sectors. Earth could be paired with a "neighbor" from Alpha Centauri, which is *light-years* away.

The solution Tradewars came up with was to have a **Cluster** that contains nearby Sectors. You can think of a cluster as a solar system or a nebula or some other celestial body.

We can update the algorithm to look something like this:

```
    create empty Universe
    fill Universe with K Sectors
    create empty Clusters
    pick random number K (where 10 > K > 0)
        add Cluster C to Clusters list
        for 0..K (where 5 > K > 0, a random number of Sectors)
            for 1..J (where 5 > J > 0, a random number of neighbors)
                pick random available Sector as N (where N != S)
                add to S.Neighbors list
            add S to C.Sectors list
            remove S from available Sector list
```

There! Now we have a selection of Clusters that each contain Sectors.

But wait, now we have a bunch of Clusters with Sectors, but how does one get from Sector 1611 to Sector 258, if they're in different clusters? The clusters have to connect somehow. 

Clearly, once we populate the Clusters, we have to link them together.

```
    for each Cluster as C
        pick random Cluster N (where N != C)
            pick random Sector S1 in C
            pick random Sector S2 in N
            add S2 to S1.Neighbors list
```

There's a minor flaw with this approach, which we can address in the next section

#### Forever Alone

In the above algorithm, you'll notice that S2 is added to S1's Neighbors list, but S1 is not added to S2's neighbor list. In other words, S1 -> S2 is a one-way ticket!

This also raises another possibility: are there any sectors that have been orphaned? Imagine the following scenario:

<div id="cluster" style="width: 100%; height: 300px;"></div>

<div id="options"></div>

See the problem? Once we go from from S1 -> S2, we can't get back to S1 -- ever! Big problem!

How do we solve this? Well, that's the $64,000 question. Tune in next time for the answer in Part II.

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