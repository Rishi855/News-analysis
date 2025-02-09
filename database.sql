
-- For creating companies
CREATE TABLE companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    groww_company_id VARCHAR(50) NOT NULL,
    api_url VARCHAR(255) NOT NULL
);

ALTER TABLE companies
ADD COLUMN search_id VARCHAR(255),
ADD COLUMN industry_code INT,
ADD COLUMN bse_script_code INT,
ADD COLUMN nse_script_code VARCHAR(50);


select * from companies ;

ALTER TABLE companies 
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;





-- for creating stock articles
-- Creating a new table to store the data except the image_url
CREATE TABLE stock_article (
    id VARCHAR(50) PRIMARY KEY,
    title TEXT NOT NULL,
    summary TEXT NOT NULL,
    url VARCHAR(255) NOT NULL,
    pub_date TIMESTAMP NOT NULL,
    source VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE stock_article
ALTER COLUMN url TYPE VARCHAR(500);  -- or larger value, depending on your needs

     
      -- Add company_id column to stock_articles
ALTER TABLE stock_article
ADD COLUMN company_id INT;

ALTER TABLE stock_article
ADD COLUMN groww_company_id TEXT;

-- Add foreign key constraint to company_id referencing companies(id)
ALTER TABLE stock_article
ADD CONSTRAINT fk_company
FOREIGN KEY (company_id) 
REFERENCES companies (id)
ON DELETE SET NULL; 
