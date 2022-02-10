--
-- PostgreSQL database dump
--

-- Dumped from database version 14.1
-- Dumped by pg_dump version 14.1

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: test; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA test;


ALTER SCHEMA test OWNER TO postgres;

--
-- Name: wallet; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA wallet;


ALTER SCHEMA wallet OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: account; Type: TABLE; Schema: wallet; Owner: postgres
--

CREATE TABLE wallet.account (
    id character varying(50) NOT NULL,
    balance numeric,
    currency character varying(50),
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE wallet.account OWNER TO postgres;

--
-- Name: payment; Type: TABLE; Schema: wallet; Owner: postgres
--

CREATE TABLE wallet.payment (
    id character varying(255) NOT NULL,
    account_id character varying(50),
    transaction_id character varying(255) NOT NULL,
    amount numeric,
    to_account character varying(50),
    from_account character varying(50),
    direction character varying(50),
    created_at timestamp without time zone NOT NULL
);


ALTER TABLE wallet.payment OWNER TO postgres;

--
-- Data for Name: account; Type: TABLE DATA; Schema: wallet; Owner: postgres
--

COPY wallet.account (id, balance, currency, created_at, updated_at) FROM stdin;
jack888	4000	USD	2022-02-09 11:32:40.417356	2022-02-09 11:32:40.417356
irin977	8000	USD	2022-02-09 11:32:40.443069	2022-02-09 11:32:40.443069
mike2167	20000	IDR	2022-02-09 11:32:40.448201	2022-02-09 11:32:40.448201
bob123	5000	USD	2022-02-09 11:19:15.988278	2022-02-09 15:21:55.396324
alice456	3000	USD	2022-02-09 11:19:40.786259	2022-02-09 15:22:01.294397
\.


--
-- Data for Name: payment; Type: TABLE DATA; Schema: wallet; Owner: postgres
--

COPY wallet.payment (id, account_id, transaction_id, amount, to_account, from_account, direction, created_at) FROM stdin;
\.


--
-- Name: account account_pk; Type: CONSTRAINT; Schema: wallet; Owner: postgres
--

ALTER TABLE ONLY wallet.account
    ADD CONSTRAINT account_pk PRIMARY KEY (id);


--
-- Name: payment payment_pk; Type: CONSTRAINT; Schema: wallet; Owner: postgres
--

ALTER TABLE ONLY wallet.payment
    ADD CONSTRAINT payment_pk PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

