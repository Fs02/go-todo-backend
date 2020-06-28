class CreateTodos < ActiveRecord::Migration[5.2]
  def change
    create_table :todos do |t|
      t.datetime :created_at
      t.datetime :updated_at
      t.string :title
      t.boolean :completed
      t.integer :order

      t.index :order
    end
  end
end
