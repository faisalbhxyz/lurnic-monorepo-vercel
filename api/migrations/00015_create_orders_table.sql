-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    orders (
        id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        student_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE,
        course_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (course_id) REFERENCES course_details (id) ON DELETE CASCADE,
        UNIQUE KEY unique_student_course (student_id, course_id),
        discount_type VARCHAR(50) NOT NULL DEFAULT 'none',
        discount DECIMAL(10, 2) NOT NULL DEFAULT 0,
        total DECIMAL(10, 2) NOT NULL,
        payment_status VARCHAR(20) NOT NULL DEFAULT 'unpaid',
        invoice_id BIGINT NOT NULL,
        payment_type VARCHAR(50) NOT NULL DEFAULT 'manual',
        customer_note TEXT,
        admin_note TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        tenant_id INT UNSIGNED NOT NULL,
        FOREIGN KEY (tenant_id) REFERENCES tenants (id) ON DELETE CASCADE,
        INDEX idx_orders_student_id (student_id),
        INDEX idx_orders_course_id (course_id),
        INDEX idx_orders_invoice_id (invoice_id),
        INDEX idx_orders_tenant_id (tenant_id)
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;

-- +goose StatementEnd