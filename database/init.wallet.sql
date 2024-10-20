/*
 Navicat Premium Dump SQL

 Source Server         : PostgreSQL14
 Source Server Type    : PostgreSQL
 Source Server Version : 140013 (140013)
 Source Catalog        : crypto
 Source Schema         : wallet

 Target Server Type    : PostgreSQL
 Target Server Version : 140013 (140013)
 File Encoding         : 65001
*/

DO
$$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_catalog.pg_namespace
        WHERE nspname = 'wallet'
    ) THEN
        EXECUTE 'CREATE SCHEMA "wallet" AUTHORIZATION "postgres";';

        -- ----------------------------
        -- Sequence structure for transation_id_seq
        -- ----------------------------
        EXECUTE 'DROP SEQUENCE IF EXISTS "wallet"."transation_id_seq";';
        EXECUTE 'CREATE SEQUENCE "wallet"."transation_id_seq" 
        INCREMENT 1
        MINVALUE  1
        MAXVALUE 2147483647
        START 1
        CACHE 1;';

        -- ----------------------------
        -- Sequence structure for user_id_seq
        -- ----------------------------
        EXECUTE 'DROP SEQUENCE IF EXISTS "wallet"."user_id_seq";';
        EXECUTE 'CREATE SEQUENCE "wallet"."user_id_seq" 
        INCREMENT 1
        MINVALUE  1
        MAXVALUE 2147483647
        START 1
        CACHE 1;';

        -- ----------------------------
        -- Table structure for transation
        -- ----------------------------
        EXECUTE 'DROP TABLE IF EXISTS "wallet"."transation";';
        EXECUTE 'CREATE TABLE "wallet"."transation" (
          "id" int4 NOT NULL GENERATED ALWAYS AS IDENTITY (
        INCREMENT 1
        MINVALUE  1
        MAXVALUE 2147483647
        START 1
        CACHE 1
        ),
          "user_id" int4 NOT NULL,
          "receiver_id" int4 NOT NULL DEFAULT 0,
          "action" varchar(20) COLLATE "pg_catalog"."default" NOT NULL,
          "money" int4 NOT NULL,
          "balance" int4 NOT NULL,
          "create_time" int4 NOT NULL
        )
        ;
        COMMENT ON COLUMN "wallet"."transation"."action" IS ''Desposit/Withdraw/Send/Receive'';';



        -- ----------------------------
        -- Table structure for user
        -- ----------------------------
        EXECUTE 'DROP TABLE IF EXISTS "wallet"."user";';
        EXECUTE 'CREATE TABLE "wallet"."user" (
          "id" int4 NOT NULL GENERATED ALWAYS AS IDENTITY (
        INCREMENT 1
        MINVALUE  1
        MAXVALUE 2147483647
        START 1
        CACHE 1
        ),
          "name" varchar(50) COLLATE "pg_catalog"."default" NOT NULL,
          "balance" int4 NOT NULL
        )
        ;';


        -- ----------------------------
        -- Alter sequences owned by
        -- ----------------------------
        EXECUTE 'ALTER SEQUENCE "wallet"."transation_id_seq" OWNED BY "wallet"."transation"."id";';
        EXECUTE 'SELECT setval(''wallet.transation_id_seq'', 655, true);';

        -- ----------------------------
        -- Alter sequences owned by
        -- ----------------------------
        EXECUTE 'ALTER SEQUENCE "wallet"."user_id_seq" OWNED BY "wallet"."user"."id";';
        EXECUTE 'SELECT setval(''wallet.user_id_seq'', 318, true);';

        -- ----------------------------
        -- Auto increment value for transation
        -- ----------------------------
        EXECUTE 'SELECT setval(''wallet.transation_id_seq'', 655, true);';

        -- ----------------------------
        -- Indexes structure for table transation
        -- ----------------------------
        EXECUTE 'CREATE INDEX "idx_action" ON "wallet"."transation" USING btree (
          "action" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
        );';
        EXECUTE 'CREATE INDEX "idx_user_id" ON "wallet"."transation" USING btree (
          "user_id" "pg_catalog"."int4_ops" ASC NULLS LAST
        );';

        -- ----------------------------
        -- Primary Key structure for table transation
        -- ----------------------------
        EXECUTE 'ALTER TABLE "wallet"."transation" ADD CONSTRAINT "transation_pkey" PRIMARY KEY ("id");';

        -- ----------------------------
        -- Auto increment value for user
        -- ----------------------------
        EXECUTE 'SELECT setval(''wallet.user_id_seq'', 318, true);';

        -- ----------------------------
        -- Uniques structure for table user
        -- ----------------------------
        EXECUTE 'ALTER TABLE "wallet"."user" ADD CONSTRAINT "uqi_name" UNIQUE ("name");';

        -- ----------------------------
        -- Primary Key structure for table user
        -- ----------------------------
        EXECUTE 'ALTER TABLE "wallet"."user" ADD CONSTRAINT "user_pkey" PRIMARY KEY ("id");';

    END IF;
END
$$;
