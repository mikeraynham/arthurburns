module Jekyll
  class GallerySection < Page
    def initialize(site, base, dir, section)
      @site = site
      @base = base
      @dir = dir
      @name = 'index.html'

      self.process(@name)
      self.read_yaml(File.join(base, '_layouts'), 'gallery_section.html')
      self.data['section'] = section
    end
  end

  class GalleryImage < Page
    def initialize(site, base, dir, section, image)
      @site = site
      @base = base
      @dir = dir
      @name = 'index.html'

      self.process(@name)
      self.read_yaml(File.join(base, '_layouts'), 'gallery_image.html')
      self.data['section'] = section
      self.data['image'] = image
    end
  end

  class GalleryPageGenerator < Generator
    safe true
    def generate(site)
      if site.layouts.key? 'gallery_section'
        dir = 'galleries'
        gallery = site.data['gallery']
        gallery.each do |section|
          site.pages << GallerySection.new(
            site,
            site.source,
            File.join(dir, section['SectionDir']),
            section
          )
          section['Images'].each do |image|
            site.pages << GalleryImage.new(
              site,
              site.source,
              File.join(dir, image['Page']),
              section,
              image
            )
          end
        end
      end
    end
  end
end
