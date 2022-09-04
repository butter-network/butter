import networkx as nx
import matplotlib.pyplot as plt

G = nx.Graph()
G.add_edges_from([(1, 2), (2, 3), (1, 3), (1, 4)])
nx.draw(G)
plt.show()
