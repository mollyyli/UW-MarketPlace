import React from "react";
import { Button, Form, FormGroup, Input, Label } from "reactstrap";

const EditListingForm = ({
  title,
  description,
  condition,
  price,
  contact,
  location,
  onHandleInputChange,
  onHandleSubmit
}) => (
    <Form onSubmit={onHandleSubmit}>
      <FormGroup>
        <Label for="title">Title</Label>
        <Input name="title" value={title} onChange={onHandleInputChange} />
      </FormGroup>
      <FormGroup>
        <Label for="description">Description</Label>
        <Input name="description" value={description} onChange={onHandleInputChange} />
      </FormGroup>
      <FormGroup>
        <Label for="condition">Condition</Label>
        <Input name="condition" value={condition} onChange={onHandleInputChange} />
      </FormGroup>
      <FormGroup>
        <Label for="price">Price</Label>
        <Input name="price" value={price} onChange={onHandleInputChange} />
      </FormGroup>
      <FormGroup>
        <Label for="contact">Contact</Label>
        <Input name="contact" value={contact} onChange={onHandleInputChange} />
      </FormGroup>
      <FormGroup>
        <Label for="location">Location</Label>
        <Input name="location" value={location} onChange={onHandleInputChange} />
      </FormGroup>
      <Button type="submit" disabled={!title || !description || !condition || !price || !contact || !location}>Edit Listing</Button>
    </Form>
  );

export default EditListingForm;
