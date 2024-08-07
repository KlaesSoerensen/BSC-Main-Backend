namespace ASP.NETCoreWebAPI.Models
{
    public class CollectionEntry
    {
        public int Id { get; set; } // Primary Key
        public required Transform Transform { get; set; } // Foreign Key (Transform) - Required, represents the position and scale within the collection
        public required GraphicalAsset GraphicalAsset { get; set; } // Foreign Key (GraphicalAsset) - Required, links to the graphical representation

        public CollectionEntry()
        {
            // Default transform values, assuming scale=1,1; zIndex=1; xOffset=0; yOffset=0 for a single item collection
            Transform = new Transform
            {
                xScale = 1,
                yScale = 1,
                zIndex = 1,
                xOffset = 0,
                yOffset = 0
            };
        }
    }
}
