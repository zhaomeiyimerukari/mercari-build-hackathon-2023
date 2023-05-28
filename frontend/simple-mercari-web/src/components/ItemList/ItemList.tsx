import React, { useState } from "react";
import { Item } from "../Item";
import { GoSearch } from 'react-icons/go';
import axios from 'axios';


interface Item {
  id: number;
  name: string;
  price: number;
  category_name: string;
}

interface Prop {
  items: Item[];
}

export const ItemList: React.FC<Prop> = (props) => {
  const [searchTerm, setSearchTerm] = useState("");

  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    const searchQuery = encodeURIComponent(e.target.value);
    const newUrl = searchQuery ? `search?name=${searchQuery}` : "/";

    // Updating the URL
    window.history.pushState(null, "", newUrl);
    setSearchTerm(e.target.value);
  };

  const filteredItems = props.items.filter((item) =>
    item.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="container-fluid">
      <form className="d-flex">
        <div className="input-group">
          <input
            className="form-control"
            type="text"
            placeholder="Search Items"
            value={searchTerm}
            onChange={handleSearch}
          />
          <span className="input-group-text">
            <GoSearch />
          </span>
        </div>
      </form>
      {filteredItems.length > 0 ? (
        filteredItems.map((item) => <Item key={item.id} item={item} />)
      ) : (
        <p>No items found.</p>
      )}
    </div>
  );
};
