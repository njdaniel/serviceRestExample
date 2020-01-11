-- 	list := []Item{
-- 		{Name: "backpack", Cost: "50ss", Description: "Carries 5 bulk"},
-- 		{Name: "sack", Cost: "80bp", Description: "Carries 5 bulk, requires at least one hand"},
-- 		{Name: "satchel", Cost: "30ss", Description: "Carries 2 bulk"},
-- 		{Name: "compass", Cost: "2gc", Description: "shows north"},
-- 		{Name: "50ft of rope", Cost: "20ss", Description: "hemp rope"},
-- 		{Name: "spyglass", Cost: "5gc", Description: "4x magnification"},
-- 		{Name: "crowbar", Cost: "18ss", Description: "tool to pry open"},
-- 		{Name: "hammer", Cost: "15ss", Description: "tool to hammer"},
-- 		{Name: "hourglass", Cost: "1gc", Description: "takes 1 hour for sand"},
-- 		{Name: "fishing rod", Cost: "10ss", Description: "tool for fishing"},
-- 		{Name: "1lb fish bait", Cost: "5bp", Description: "bugs to help catch fish"},
-- 	}


INSERT INTO items (
    name, cost, description
    ) VALUES ('backpack', '50ss', 'Carries 5 bulk'),
    ('satchel', '30ss', 'Carries 2 bulk');