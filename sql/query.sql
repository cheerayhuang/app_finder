\o /tmp/out.txt;
select distinct sdkbox_app_package_id
from event_log
where timestamp > extract(epoch from '2016-12-13'::date) and timestamp < extract(epoch from '2016-12-14'::date)
limit 100;
