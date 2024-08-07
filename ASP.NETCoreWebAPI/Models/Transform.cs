namespace ASP.NETCoreWebAPI.Models
{
    public class Transform
    {
        public int Id { get; set; } // Primary Key
        public int zIndex { get; set; } // Attribute - Layering index, where 0<n is for background and n>1000 is for UI
        public float xOffset { get; set; } // Attribute - X-axis offset from colony "origo"
        public float yOffset { get; set; } // Attribute - Y-axis offset from colony "origo"
        public float xScale { get; set; } // Attribute - Scale factor on the X-axis
        public float yScale { get; set; } // Attribute - Scale factor on the Y-axis
    }
}
