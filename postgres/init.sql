
CREATE TABLE IF NOT EXISTS subscriptions (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL,
    price INT NOT NULL CHECK (price >= 0),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             );

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_service
    ON subscriptions (user_id, service_name);

CREATE INDEX IF NOT EXISTS idx_subscriptions_dates
    ON subscriptions (start_date, end_date);