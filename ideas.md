a lot of the ideas for know host list maintatnance and information retrieval were pulled from GOFAI techniques I studied in CS255 (Nathan's AI module) - I found as I was developping the system I kept borrowing ideas from the field of AI

**know host list management** 

Problem: optimal list of known hosts to maintain information availability, fast information retrieval while not losing sight of smaller less important nodes (you don't want to push nodes out of the network i.e. node to be forgotten) - so the idea is to have a diverse list of known hosts (both high and low uptime and lost of available storage i.e. will be able to query node to store something + little storage availability i.e. has lost of information so higher probability of having the information we want)

- state of a given node host list
- default behaviour is to use as much memory as possible to store known host - might as well have a complete graph if we can
- if full
  - creating permutations of the list (when trying to add a new known host)
  - see what combination of hosts maximises the list quality value function
    - the value function take into consideration the remote host uptime, available memory, nb of known hosts

information retrieval

- naive BFS
- directed BFS - graph search problem where we have a valu function i.e. determine which node has the highest probability of either having the information of knowing someone that does



Presentation structure

- present the project (like a product reveal) - introduce the framework a user perspective
  - How is it a framework i.e. becomes the user inputs code into the framework (while in a library you simply call the library)
  - Using the framework
- technical - the key problems when designing decentralised unstructiured peer-to-peer architecture systems and how the framework addresses them
  - Dicovery - cold start problem
  - known host management - heavily AI inspired (GOFAI)
  - overlay network for persistent information 
    - information retrieval - heavily AI inspired (GOFAI)
    - persistent fault-tolerant storage (high availability)
  - NAT traversal