import axios from 'axios';

export const get = async (url, params, token) => {
    // Now token is passed as a parameter
    const result = await axios.get(url, {
      headers: {
        Authorization: `Bearer ${token}`
      },
      params: params // Add any existing parameters
    });
    return result.data;
  };  

export const post = async (url, params) => {
    const token = localStorage.getItem('token');

    const result = await axios.post(url, params, {
        headers: {
            Authorization: `Bearer ${token}`
        }
    });
    return result.data;
};
