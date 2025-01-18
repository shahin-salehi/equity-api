/*
    (1) Check what IDs is present in db but not in list
    (2) take (1) and set removed field date
    (3) return to caller the IDs present in list but not in db. 

   Notes for me name of the variable in return table doesnt matter its just a type
   internally its called that but to caller they get whatever they run in their select... 
 */

CREATE OR REPLACE FUNCTION listings_delta(ids text[]) RETURNS TABLE(output_list text[]) AS $$

DECLARE
   input_ids ALIAS for $1;
BEGIN

    -- One line update IDs that are present in db but not in scraped to removed. 
    UPDATE listings SET removed = NOW() WHERE listing_id IN (SELECT listing_id FROM listings WHERE NOT (listing_id = ANY(input_ids))) AND removed=NULL;
   
    -- return ids exlusive to scraped ids
	RETURN QUERY SELECT ARRAY((SELECT * FROM unnest(input_ids) EXCEPT ALL SELECT listing_id FROM listings))::text[];

END;

$$ LANGUAGE plpgsql;