DROP TABLE "users";

ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT "accounts_owner_fkey";
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT "owner_currency_key";

