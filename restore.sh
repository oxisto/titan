USER=postgres
VERSION=`cat sde.version`
HOST=$1

# EVE SDE
psql -U $USER -h $HOST titan -c 'DROP SCHEMA evesde CASCADE; CREATE SCHEMA evesde'
pg_restore -U $USER -h $HOST \
-t invTypes \
-t invGroups \
-t invCategories \
-t invMetaTypes \
-t industryBlueprints \
-t industryActivity \
-t industryActivityProducts \
-t industryActivitySkills \
-t industryActivityMaterials \
-d titan sde-$VERSION

# Our tables
psql -U $USER -h $HOST titan < sql/public.sql

