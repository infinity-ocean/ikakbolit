INSERT INTO developers (id, name, department, geolocation, last_known_ip, is_available)
SELECT
  uuid_generate_v4(),
  (ARRAY['James','Mary','John','Patricia','Robert'])[floor(random()*5+1)] 
    || ' ' || 
  (ARRAY['Smith','Johnson','Williams','Brown','Jones'])[floor(random()*5+1)],
  (ARRAY['backend','frontend','ios','android'])[floor(random()*4+1)],
  ST_MakePoint(random()*360-180, random()*180-90)::geography,
  (
    floor(random()*256)::int || '.' ||
    floor(random()*256)::int || '.' ||
    floor(random()*256)::int || '.' ||
    floor(random()*256)::int
  )::inet,
  (random() < 0.5)
FROM generate_series(1,1000);