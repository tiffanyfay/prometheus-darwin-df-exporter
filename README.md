# prometheus-darwin-df-exporter

This came about because the [Prometheus Node Exporter](https://github.com/prometheus/node_exporter/tree/master) does not show filesystem used storage. It uses exec to run `df -k` since `unix.Statfs` also doesn't provide this.

Here's what `df -k` gives on my machine:
```
Filesystem     1024-blocks      Used Available Capacity iused      ifree %iused  Mounted on
/dev/disk1s5s1   488245288  15057096 228509680     7%  502388 2285096800    0%   /
devfs                  202       202         0   100%     702          0  100%   /dev
/dev/disk1s4     488245288  11535436 228509680     5%      12 2285096800    0%   /System/Volumes/VM
/dev/disk1s2     488245288    380564 228509680     1%    1992 2285096800    0%   /System/Volumes/Preboot
/dev/disk1s6     488245288      3116 228509680     1%      17 2285096800    0%   /System/Volumes/Update
/dev/disk1s1     488245288 231508056 228509680    51% 1508773 2285096800    0%   /System/Volumes/Data
map auto_home            0         0         0   100%       0          0  100%   /System/Volumes/Data/home
```