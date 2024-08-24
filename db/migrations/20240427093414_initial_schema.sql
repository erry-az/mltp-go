-- Create enum type "transaction_type"
CREATE TYPE "transaction_type" AS ENUM ('credit', 'debit');
-- Create enum type "transaction_name"
CREATE TYPE "transaction_name" AS ENUM ('transfer', 'top_up');
-- Create "users" table
CREATE TABLE "users" ("id" serial NOT NULL, "username" character varying(255) NOT NULL, "fullname" character varying(255) NOT NULL, "balance" bigint NOT NULL DEFAULT 0, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NULL, PRIMARY KEY ("id"), CONSTRAINT "users_username_key" UNIQUE ("username"));
-- Create "transactions" table
CREATE TABLE "transactions" ("id" serial NOT NULL, "user_id" integer NOT NULL, "amount" bigint NOT NULL, "name" "transaction_name" NOT NULL, "type" "transaction_type" NOT NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("id"), CONSTRAINT "transactions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
