#!/bin/bash
ESHOST=http://127.0.0.1:9200
ESDIR="/web/elasticsearch-0.90.9"
INDEXNAME="scripts" # this had better match the index name in ES
INDEXDIR="$ESDIR/data/elasticsearch/nodes/0/indices/"
BACKUPDIR="$ESDIR/backups/"
YEARMONTH=`date +"%Y-%m"`

# create mapping file with index settings. this metadata is required by ES to use index file data
echo -n "Backing up metadata... "
curl -XGET -o /tmp/mapping "$ESHOST/$INDEXNAME/_mapping?pretty=true" > /dev/null 2>&1
sed -i '1,2d' /tmp/mapping #strip the first two lines of the metadata
echo '{"settings":{"number_of_shards":5,"number_of_replicas":1},"mappings":{' >> /tmp/mappost 
# prepend hardcoded settings metadata to index-specific metadata
cat /tmp/mapping >> /tmp/mappost
echo "DONE!"

# now lets tar up our data files. these are huge, so lets be nice
echo -n "Backing up data files (this may take some time)... "
mkdir -p $BACKUPDIR
cd $INDEXDIR
nice -n 19 tar czf $BACKUPDIR/$INDEXNAME.tar.gz $INDEXNAME 
echo "DONE!"

echo -n "Creating restore script... "
# time to create our restore script! oh god scripts creating scripts, this never ends well...
cat << EOF >> $BACKUPDIR/$INDEXNAME-restore.sh
#!/bin/bash
# this script requires $INDEXNAME.tar.gz and will restore it into elasticsearch
# it is ESSENTIAL that the index you are restoring does NOT exist in ES. delete it
# if it does BEFORE trying to restore data.

# create index and mapping
echo -n "Creating index and mappings... "
curl -XPUT 'http://localhost:9200/$INDEXNAME/' -d '`cat /tmp/mappost`' > /dev/null 2>&1
echo "DONE!"

# extract our data files into place
echo -n "Restoring index (this may take a while)... "
cd $INDEXDIR
tar xzf $BACKUPDIR/$INDEXNAME.tar.gz
echo "DONE!"

# restart ES to allow it to open the new dir and file data
echo -n "Restarting Elasticsearch... "
/etc/init.d/es restart
echo "DONE!"
EOF
echo "DONE!" # restore script done

# cleanup tmp files
rm /tmp/mappost
rm /tmp/mapping
