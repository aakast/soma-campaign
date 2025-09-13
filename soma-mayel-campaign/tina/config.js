import { defineConfig } from 'tinacms';

export default defineConfig({
  branch: process.env.TINA_PUBLIC_BRANCH || 'main',
  clientId: process.env.TINA_PUBLIC_CLIENT_ID,
  token: process.env.TINA_PUBLIC_TOKEN,
  
  build: {
    outputFolder: 'admin',
    publicFolder: 'static',
  },
  
  media: {
    tina: {
      mediaRoot: 'images',
      publicFolder: 'static',
    },
  },
  
  schema: {
    collections: [
      {
        name: 'post',
        label: 'Blog Posts',
        path: 'content/posts',
        format: 'json',
        fields: [
          {
            type: 'string',
            name: 'title',
            label: 'Title',
            isTitle: true,
            required: true,
          },
          {
            type: 'string',
            name: 'slug',
            label: 'Slug',
            required: true,
          },
          {
            type: 'datetime',
            name: 'date',
            label: 'Date',
            required: true,
          },
          {
            type: 'string',
            name: 'author',
            label: 'Author',
            options: ['Soma Mayel', 'Radikale Venstre Fredensborg'],
          },
          {
            type: 'string',
            name: 'excerpt',
            label: 'Excerpt',
            ui: {
              component: 'textarea',
            },
          },
          {
            type: 'image',
            name: 'image',
            label: 'Featured Image',
          },
          {
            type: 'string',
            name: 'tags',
            label: 'Tags',
            list: true,
          },
          {
            type: 'boolean',
            name: 'isFeatured',
            label: 'Featured Post',
          },
          {
            type: 'rich-text',
            name: 'content',
            label: 'Content',
            isBody: true,
          },
        ],
      },
      {
        name: 'page',
        label: 'Pages',
        path: 'content/pages',
        format: 'json',
        fields: [
          {
            type: 'string',
            name: 'title',
            label: 'Title',
            isTitle: true,
            required: true,
          },
          {
            type: 'rich-text',
            name: 'content',
            label: 'Content',
            isBody: true,
          },
        ],
      },
      {
        name: 'settings',
        label: 'Site Settings',
        path: 'content',
        format: 'json',
        ui: {
          allowedActions: {
            create: false,
            delete: false,
          },
        },
        fields: [
          {
            type: 'string',
            name: 'siteTitle',
            label: 'Site Title',
          },
          {
            type: 'string',
            name: 'siteDescription',
            label: 'Site Description',
            ui: {
              component: 'textarea',
            },
          },
          {
            type: 'object',
            name: 'contact',
            label: 'Contact Information',
            fields: [
              {
                type: 'string',
                name: 'email',
                label: 'Email',
              },
              {
                type: 'string',
                name: 'phone',
                label: 'Phone',
              },
              {
                type: 'string',
                name: 'facebook',
                label: 'Facebook URL',
              },
              {
                type: 'string',
                name: 'address',
                label: 'Address',
              },
            ],
          },
          {
            type: 'object',
            name: 'hero',
            label: 'Hero Section',
            fields: [
              {
                type: 'string',
                name: 'title',
                label: 'Title',
              },
              {
                type: 'string',
                name: 'subtitle',
                label: 'Subtitle',
              },
              {
                type: 'image',
                name: 'videoUrl',
                label: 'Video URL',
              },
            ],
          },
        ],
      },
    ],
  },
});