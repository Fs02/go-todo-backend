class CreateScores < ActiveRecord::Migration[5.2]
  def change
    create_table :scores do |t|
      t.datetime :created_at
      t.datetime :updated_at
      t.integer :total_point
    end
  end
end
