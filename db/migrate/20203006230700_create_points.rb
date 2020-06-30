class CreatePoints < ActiveRecord::Migration[5.2]
  def change
    create_table :points do |t|
      t.datetime :created_at
      t.datetime :updated_at
      t.string :name
      t.integer :count

      t.belongs_to :score
    end
  end
end
