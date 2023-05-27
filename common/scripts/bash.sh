
# 分配文件夹
for i in `seq 1 6`
do
    rm -rf /tmp/$i/objects/*
    rm -rf /tmp/$i/temp/*
    rm -rf /tmp/$i/garbage/*
done

for i in `seq 1 6`
do
    mkdir -p /log/$i/
done

mkdir -p /log/21
mkdir -p /log/22

