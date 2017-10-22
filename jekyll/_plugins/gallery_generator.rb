module Jekyll
    class GallerySection < Page
        def initialize(site, base, dir, section, prev_section, next_section)
            @site = site
            @base = base
            @dir  = dir
            @name = 'index.html'

            self.process(@name)
            self.read_yaml(File.join(base, '_layouts'), 'gallery_section.html')
            self.data['section']      = section
            self.data['prev_section'] = prev_section if !prev_section.empty?
            self.data['next_section'] = next_section if !next_section.empty?
        end
    end

    class GalleryImage < Page
        def initialize(site, base, dir, section, image, prev_image, next_image)
            @site = site
            @base = base
            @dir  = dir
            @name = 'index.html'

            self.process(@name)
            self.read_yaml(File.join(base, '_layouts'), 'gallery_image.html')
            self.data['section']    = section
            self.data['image']      = image
            self.data['prev_image'] = prev_image if !prev_image.empty?
            self.data['next_image'] = next_image if !next_image.empty?
        end
    end

    class GalleryPageGenerator < Generator
        safe true
        def generate(site)
            if !site.layouts.key? 'gallery_section'
                return
            end

            dir          = 'galleries'
            gallery      = site.data['gallery']
            last_section = gallery.length - 1
            prev_section = {}
            prev_image   = {}

            gallery.each_index do |section_index|
                section      = gallery[section_index]
                images       = section['Images']
                last_image   = images.length - 1
                next_section = section_index < last_section ? gallery[section_index + 1] : {}

                site.pages << GallerySection.new(
                    site,
                    site.source,
                    File.join(dir, section['SectionDir']),
                    section,
                    prev_section,
                    next_section
                )

                images.each_index do |image_index|
                    image      = images[image_index]
                    next_image = {}

                    if image_index < last_image
                        next_image = images[image_index + 1]
                    elsif !next_section.empty?
                        next_image = next_section['Images'][0]
                    end

                    puts section['SectionTitle']
                    puts image['Page']
                    puts next_image['Page']

                    site.pages << GalleryImage.new(
                        site,
                        site.source,
                        File.join(dir, image['Page']),
                        section,
                        image,
                        prev_image,
                        next_image
                    )
                    prev_image = image
                end

                prev_section = section
            end
        end
    end
end
