import React, { useState } from 'react';
import axios from 'axios';

export const SearchBar = () => {
  const [searchTerm, setSearchTerm] = useState('');

  const handleSearch = async (e: any) => {
    e.preventDefault();
    if (searchTerm.trim() !== '') {
      try {
        const response = await axios.get(`/search?name=${encodeURIComponent(searchTerm)}`);
        const searchResults = response.data;
        // Process and display searchResults in your frontend as needed
      } catch (error) {
        // Handle error
        console.error('Error searching for items:', error);
      }
    }
  };

  return (
    <div className="container-fluid">
      <form className="d-flex"  onSubmit={handleSearch}>
        <input
          className="form-control me-2"
          type="search"
          placeholder="Search for items"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
        <button className="btn btn-outline-danger" type="submit">Search</button>
      </form>
    </div>
  );
};
