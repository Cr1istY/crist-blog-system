--
-- PostgreSQL database dump
--

\restrict KmhkaR4yeA4xrRzagVDWkPMv6V7BuiCz6yX7koa7Rwbxxo9EnraysFydL9aSF5A

-- Dumped from database version 13.23
-- Dumped by pg_dump version 13.23

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
-- Name: admin; Type: SCHEMA; Schema: -; Owner: crist
--

CREATE SCHEMA admin;


ALTER SCHEMA admin OWNER TO crist;

--
-- Name: blog; Type: SCHEMA; Schema: -; Owner: crist
--

CREATE SCHEMA blog;


ALTER SCHEMA blog OWNER TO crist;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: post_status_enum; Type: TYPE; Schema: public; Owner: crist
--

CREATE TYPE public.post_status_enum AS ENUM (
    'draft',
    'published',
    'archived'
);


ALTER TYPE public.post_status_enum OWNER TO crist;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: refresh_tokens; Type: TABLE; Schema: admin; Owner: crist
--

CREATE TABLE admin.refresh_tokens (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    token_hash text NOT NULL,
    user_agent text,
    ip_address inet,
    expires_at timestamp with time zone NOT NULL,
    revoked boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    province text
);


ALTER TABLE admin.refresh_tokens OWNER TO crist;

--
-- Name: users; Type: TABLE; Schema: admin; Owner: crist
--

CREATE TABLE admin.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    username text NOT NULL,
    password_hash text NOT NULL,
    nickname text,
    email text,
    avatar text,
    bio text,
    is_admin boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    CONSTRAINT users_username_check CHECK ((char_length(username) >= 3))
);


ALTER TABLE admin.users OWNER TO crist;

--
-- Name: categories; Type: TABLE; Schema: blog; Owner: crist
--

CREATE TABLE blog.categories (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    slug text NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT now(),
    parent_id uuid,
    deleted_flag boolean DEFAULT false NOT NULL,
    CONSTRAINT categories_slug_check CHECK ((slug ~ '^[a-z0-9]+(-[a-z0-9]+)*$'::text))
);


ALTER TABLE blog.categories OWNER TO crist;

--
-- Name: images; Type: TABLE; Schema: blog; Owner: crist
--

CREATE TABLE blog.images (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    url character varying(512) NOT NULL,
    filename character varying(255) NOT NULL,
    size bigint NOT NULL,
    width integer,
    height integer,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE blog.images OWNER TO crist;

--
-- Name: posts; Type: TABLE; Schema: blog; Owner: crist
--

CREATE TABLE blog.posts (
    id bigint NOT NULL,
    user_id uuid NOT NULL,
    title text NOT NULL,
    slug text NOT NULL,
    content text,
    excerpt text,
    status public.post_status_enum DEFAULT 'draft'::public.post_status_enum,
    category_id uuid,
    tags text[],
    views integer DEFAULT 0,
    likes integer DEFAULT 0,
    published_at timestamp with time zone,
    meta_title text,
    meta_description text,
    search_vector tsvector,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    thumbnail text,
    is_pinned boolean DEFAULT false,
    pinned_order integer DEFAULT 0,
    pinned_until timestamp with time zone
);


ALTER TABLE blog.posts OWNER TO crist;

--
-- Name: posts_id_seq; Type: SEQUENCE; Schema: blog; Owner: crist
--

CREATE SEQUENCE blog.posts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE blog.posts_id_seq OWNER TO crist;

--
-- Name: posts_id_seq; Type: SEQUENCE OWNED BY; Schema: blog; Owner: crist
--

ALTER SEQUENCE blog.posts_id_seq OWNED BY blog.posts.id;


--
-- Name: tags; Type: TABLE; Schema: blog; Owner: crist
--

CREATE TABLE blog.tags (
    name text NOT NULL,
    slug text NOT NULL,
    description text,
    post_count integer DEFAULT 0,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    CONSTRAINT tags_slug_check CHECK ((slug ~ '^[a-z0-9]+(-[a-z0-9]+)*$'::text))
);


ALTER TABLE blog.tags OWNER TO crist;

--
-- Name: tweet_images; Type: TABLE; Schema: blog; Owner: crist
--

CREATE TABLE blog.tweet_images (
    tweet_id uuid NOT NULL,
    image_id uuid NOT NULL,
    display_order integer DEFAULT 0 NOT NULL
);


ALTER TABLE blog.tweet_images OWNER TO crist;

--
-- Name: tweets; Type: TABLE; Schema: blog; Owner: crist
--

CREATE TABLE blog.tweets (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    content text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_flag boolean DEFAULT false NOT NULL,
    likes integer
);


ALTER TABLE blog.tweets OWNER TO crist;

--
-- Name: posts id; Type: DEFAULT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.posts ALTER COLUMN id SET DEFAULT nextval('blog.posts_id_seq'::regclass);


--
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: admin; Owner: crist
--

ALTER TABLE ONLY admin.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: admin; Owner: crist
--

ALTER TABLE ONLY admin.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: admin; Owner: crist
--

ALTER TABLE ONLY admin.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: admin; Owner: crist
--

ALTER TABLE ONLY admin.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: categories categories_name_key; Type: CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.categories
    ADD CONSTRAINT categories_name_key UNIQUE (name);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: categories categories_slug_key; Type: CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.categories
    ADD CONSTRAINT categories_slug_key UNIQUE (slug);


--
-- Name: images images_pkey; Type: CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.images
    ADD CONSTRAINT images_pkey PRIMARY KEY (id);


--
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- Name: posts posts_slug_key; Type: CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.posts
    ADD CONSTRAINT posts_slug_key UNIQUE (slug);


--
-- Name: tags tags_pkey; Type: CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (name);


--
-- Name: tags tags_slug_key; Type: CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.tags
    ADD CONSTRAINT tags_slug_key UNIQUE (slug);


--
-- Name: tweet_images tweet_images_pkey; Type: CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.tweet_images
    ADD CONSTRAINT tweet_images_pkey PRIMARY KEY (tweet_id, image_id);


--
-- Name: tweets tweets_pkey; Type: CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.tweets
    ADD CONSTRAINT tweets_pkey PRIMARY KEY (id);


--
-- Name: idx_refresh_tokens_expires; Type: INDEX; Schema: admin; Owner: crist
--

CREATE INDEX idx_refresh_tokens_expires ON admin.refresh_tokens USING btree (expires_at);


--
-- Name: idx_refresh_tokens_token_hash; Type: INDEX; Schema: admin; Owner: crist
--

CREATE INDEX idx_refresh_tokens_token_hash ON admin.refresh_tokens USING btree (token_hash);


--
-- Name: idx_refresh_tokens_user; Type: INDEX; Schema: admin; Owner: crist
--

CREATE INDEX idx_refresh_tokens_user ON admin.refresh_tokens USING btree (user_id);


--
-- Name: idx_posts_active_pinned; Type: INDEX; Schema: blog; Owner: crist
--

CREATE INDEX idx_posts_active_pinned ON blog.posts USING btree (is_pinned DESC, pinned_order, published_at DESC) WHERE ((deleted_at IS NULL) AND (status = 'published'::public.post_status_enum) AND (is_pinned = true));


--
-- Name: idx_posts_pinned_until; Type: INDEX; Schema: blog; Owner: crist
--

CREATE INDEX idx_posts_pinned_until ON blog.posts USING btree (pinned_until) WHERE ((pinned_until IS NOT NULL) AND (is_pinned = true));


--
-- Name: idx_tweet_images_tweet_id; Type: INDEX; Schema: blog; Owner: crist
--

CREATE INDEX idx_tweet_images_tweet_id ON blog.tweet_images USING btree (tweet_id);


--
-- Name: idx_tweets_created_at; Type: INDEX; Schema: blog; Owner: crist
--

CREATE INDEX idx_tweets_created_at ON blog.tweets USING btree (created_at DESC);


--
-- Name: idx_tweets_user_id; Type: INDEX; Schema: blog; Owner: crist
--

CREATE INDEX idx_tweets_user_id ON blog.tweets USING btree (user_id);


--
-- Name: refresh_tokens refresh_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: admin; Owner: crist
--

ALTER TABLE ONLY admin.refresh_tokens
    ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES admin.users(id) ON DELETE CASCADE;


--
-- Name: categories categories_parent_id_fkey; Type: FK CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.categories
    ADD CONSTRAINT categories_parent_id_fkey FOREIGN KEY (parent_id) REFERENCES blog.categories(id) ON DELETE CASCADE;


--
-- Name: posts posts_category_id_fkey; Type: FK CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.posts
    ADD CONSTRAINT posts_category_id_fkey FOREIGN KEY (category_id) REFERENCES blog.categories(id) ON DELETE SET NULL;


--
-- Name: posts posts_user_id_fkey; Type: FK CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.posts
    ADD CONSTRAINT posts_user_id_fkey FOREIGN KEY (user_id) REFERENCES admin.users(id) ON DELETE CASCADE;


--
-- Name: tweet_images tweet_images_image_id_fkey; Type: FK CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.tweet_images
    ADD CONSTRAINT tweet_images_image_id_fkey FOREIGN KEY (image_id) REFERENCES blog.images(id) ON DELETE CASCADE;


--
-- Name: tweet_images tweet_images_tweet_id_fkey; Type: FK CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.tweet_images
    ADD CONSTRAINT tweet_images_tweet_id_fkey FOREIGN KEY (tweet_id) REFERENCES blog.tweets(id) ON DELETE CASCADE;


--
-- Name: tweets tweets_user_id_fkey; Type: FK CONSTRAINT; Schema: blog; Owner: crist
--

ALTER TABLE ONLY blog.tweets
    ADD CONSTRAINT tweets_user_id_fkey FOREIGN KEY (user_id) REFERENCES admin.users(id);


--
-- PostgreSQL database dump complete
--

\unrestrict KmhkaR4yeA4xrRzagVDWkPMv6V7BuiCz6yX7koa7Rwbxxo9EnraysFydL9aSF5A

