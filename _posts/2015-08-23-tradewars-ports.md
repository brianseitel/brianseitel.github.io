---
layout: post
title: Remaking Tradewars - Part V
subtitle: Outposts and Traders
date: 2015-08-23 09:42:52
disqus: n
---

### Tradewars -- Trading what?

Half the *point* of Tradewars is the act of trading! Now that we've established a universe, we need to populate it with two things: ports and traders!

In Tradewars, each port sells three things: Fuel, Ore, and Equipment. The port also gets designated a Class, and the Class identifies which type of trades they make. The description is delineated as [BBS]. The **B** stands for Buy, and the **S** stands for Sell. The letters occur in the order of Fuel, Ore, and Equipment.

For instance, [BBS] means it buys fuel, buys ore, and sells equipment.

If you visit a port that sells Equipment, you cannot *buy* equipment from that port.

Enough about that. Here's a list of the possible Classes:

<table class="table">
    <thead>
        <th>Class</th>
        <th>Description</th>
        <th>Fuel</th>
        <th>Ore</th>
        <th>Equipment</th>
    </thead>
    <tbody>
        <tr>
            <td>1</td>
            <td>BBS</td>
            <td>Buy</td>
            <td>Buy</td>
            <td>Sell</td>
        </tr>
        <tr>
            <td>2</td>
            <td>BSB</td>
            <td>Buy</td>
            <td>Sell</td>
            <td>Buy</td>
        </tr>
        <tr>
            <td>3</td>
            <td>SBB</td>
            <td>Sell</td>
            <td>Buy</td>
            <td>Buy</td>
        </tr>
        <tr>
            <td>4</td>
            <td>SSB</td>
            <td>Sell</td>
            <td>Sell</td>
            <td>Buy</td>
        </tr>
        <tr>
            <td>5</td>
            <td>SBS</td>
            <td>Sell</td>
            <td>Buy</td>
            <td>Sell</td>
        </tr>
        <tr>
            <td>6</td>
            <td>BSS</td>
            <td>Buy</td>
            <td>Sell</td>
            <td>Sell</td>
        </tr>
        <tr>
            <td>7</td>
            <td>SSS</td>
            <td>Sell</td>
            <td>Sell</td>
            <td>Sell</td>
        </tr>
        <tr>
            <td>8</td>
            <td>BBB</td>
            <td>Buy</td>
            <td>Buy</td>
            <td>Buy</td>
        </tr>
    </tbody>
</table>

Make sense? Great.

### Port Theory

Like I said earlier, half the point of the game is to trade. Since some ports only buy certain goods and other ports only sell certain goods, it stands to reason that a savvy trader might find two ports near each other with reasonable prices and travel back and forth between the two of them.

If we have ports too close together, the game gets too easy. If they're too far apart, it takes too long to play.

In other words, we need to be careful about exactly how many outposts we set up.

The second concern is pricing. If Port A sells Fuel for 100 credits, and Port B buys Fuel for 200 credits, then it's profitable to travel to Port A, buy up a bunch of Fuel, then go to Port B and sell it all.

On teh other hand, if Port A sells Fuel for 200 credits, and port B buys Fuel for 100 credits, then it's not worth it at all. You'd lose 100 credits per Fuel!

On yet another hand (apparently we have three hands), you don't want it to *always* be profitable. That would be too easy. We'll need to add an element of randomness.

### Generating Ports

When we generate ports, we have to think of a few numbers.

#### Port Density and Distribution

Let's say we initialize our universe with 100 sectors. (Clusters don't really matter -- if we choose random sectors, our ports should be randomly distributed.)

With 100 sectors and a port density of about 40%, that gives us 40 sectors. Sounds about right.

Not all ports are created equal. The original Tradewars created Ports in roughly this distribution:

<table class="table">
    <thead>
        <th>Class</th>
        <th>Description</th>
        <th>Chance</th>
    </thead>
    <tbody>
        <tr>
            <td>1</td>
            <td>BBS</td>
            <td>20%</td>
        </tr>
        <tr>
            <td>2</td>
            <td>BSB</td>
            <td>20%</td>
        </tr>
        <tr>
            <td>3</td>
            <td>SBB</td>
            <td>20%</td>
        </tr>
        <tr>
            <td>4</td>
            <td>SSB</td>
            <td>10%</td>
        </tr>
        <tr>
            <td>5</td>
            <td>SBS</td>
            <td>10%</td>
        </tr>
        <tr>
            <td>6</td>
            <td>BSS</td>
            <td>10%</td>
        </tr>
        <tr>
            <td>7</td>
            <td>SSS</td>
            <td>5%</td>
        </tr>
        <tr>
            <td>8</td>
            <td>BBB</td>
            <td>5%</td>
        </tr>
    </tbody>
</table>

#### Prices

The next thing we need to decide is the pricing structure. In Tradewars, Fuel is cheapest and Equipment is the most expensive, with Ore falling somewhere in between.

One strategy is to do something like: `RAND(STANDARD) + MIN`. You pick a "standard" price and then specify a minimum. Pick a random number between 0 and `STANDARD`, then add the minimum. This gives you a final number between `MIN` and `STANDARD + MIN`.

Here are the values I picked, which may or may not line up to Tradewars standards:

<table class="table">
    <thead>
        <th>Class</th>
        <th>MIN</th>
        <th>STANDARD</th>
        <th>Lowest</th>
        <th>Highest</th>
    </thead>
    <tbody>
        <tr>
            <td>Fuel</td>
            <td>50</td>
            <td>100</td>
            <td>50</td>
            <td>150</td>
        </tr>
        <tr>
            <td>Ore</td>
            <td>150</td>
            <td>500</td>
            <td>150</td>
            <td>650</td>
        </tr>
        <tr>
            <td>Equipment</td>
            <td>750</td>
            <td>2500</td>
            <td>750</td>
            <td>3250</td>
        </tr>
    </tbody>
</table>

#### Inventory

The final thing we need to consider about a port is its inventory: how many of each item do we have?

Let's go back to our original thought: if Port A sold Fuel for 100 and Port B bought Fuel for 200, then the trader could go back and forth and make a huge profit. This is good in the short term, but it makes the game way too easy if you can just do it forever! The solution: limited inventory.

We can use exactly the same algorithm as above to generate random inventories for a given Port.

For prices, Fuel was the cheapest and Equipment was the most expensive. The basic economic theory of supply and demand tells us that Fuel, being the cheapest, would be the most plentiful, while Equipment, being the most expensive, would be relatively rare.

For our inventory, here are the values I picked, which, again, may or may not line up to Tradewars standards:

<table class="table">
    <thead>
        <th>Class</th>
        <th>MIN</th>
        <th>STANDARD</th>
        <th>Lowest</th>
        <th>Highest</th>
    </thead>
    <tbody>
        <tr>
            <td>Fuel</td>
            <td>50</td>
            <td>100</td>
            <td>50</td>
            <td>150</td>
        </tr>
        <tr>
            <td>Ore</td>
            <td>10</td>
            <td>50</td>
            <td>10</td>
            <td>60</td>
        </tr>
        <tr>
            <td>Equipment</td>
            <td>0</td>
            <td>25</td>
            <td>0</td>
            <td>25</td>
        </tr>
    </tbody>
</table>

You'll notice that it's theoretically possible for a Port to not have any Equipment at all! That sucks, but fortunately, it'll only be temporary. All ports have to restock, eventually!

### The Algorithm

This one is pretty straightforward:

```
TOTAL_SECTORS = 100
PORT_DENSITY  = 0.4
NUM_PORTS     = TOTAL_SECTORS * PORT_DENSITY

for 0..NUM_PORTS
    port = new Port
    chance = Random number between 0 and 100
    switch(true) {
        case chance < 20:  port.class = 1; break; // 20% chance of Class 1
        case chance < 40:  port.class = 2; break; // 20% chance of Class 2
        case chance < 60:  port.class = 3; break; // 20% chance of Class 3
        case chance < 70:  port.class = 4; break; // 10% chance of Class 4
        case chance < 80:  port.class = 5; break; // 10% chance of Class 5
        case chance < 90:  port.class = 6; break; // 10% chance of Class 6
        case chance < 95:  port.class = 7; break; // 5% chance of Class 7
        case chance < 100: port.class = 8; break; // 5% chance of Class 8
    }

    port.price.fuel      = Random(100) + 50;
    port.price.ore       = Random(500) + 150;
    port.price.equipment = Random(2500) + 750;

    port.inventory.fuel      = Random(100) + 50;
    port.inventory.ore       = Random(50) + 10;
    port.inventory.equipment = Random(25) + 0;

    add port to random sector in Universe
end
```

That should be it!